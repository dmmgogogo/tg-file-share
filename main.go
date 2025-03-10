package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"tg-file-share/bot"
	"tg-file-share/conf"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// Telegram API 响应结构体
type GetFileResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		FilePath string `json:"file_path"`
	} `json:"result"`
}

// FileController 处理文件相关的请求
type FileController struct {
	web.Controller
}

// 获取 Telegram 文件路径
func GetFilePath(fileID string) (string, error) {
	apiURL := fmt.Sprintf("%s/bot%s/getFile", conf.TelegramAPI, conf.TelegramBotToken)
	resp, err := http.PostForm(apiURL, map[string][]string{"file_id": {fileID}})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result GetFileResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if !result.Ok {
		return "", fmt.Errorf("failed to get file path")
	}

	return result.Result.FilePath, nil
}

// Get 处理文件访问请求
func (c *FileController) Get() {
	fileID := c.Ctx.Input.Param(":file_id")
	if fileID == "" {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.Body([]byte("file_id 参数不能为空"))
		return
	}

	// 获取 file_path
	filePath, err := GetFilePath(fileID)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.Body([]byte("无法获取文件路径: " + err.Error()))
		return
	}

	logs.Debug("fileID: %s, filePath: %s", fileID, filePath)

	// 构造 Telegram 访问 URL
	// fileURL := fmt.Sprintf("%s/file/bot%s/%s", conf.TelegramAPI, conf.TelegramBotToken, filePath)
	fileURL := fmt.Sprintf("%s/d/%s", conf.FileServerURL, filePath)

	// 301 重定向
	c.Redirect(fileURL, 301)
	// c.Ctx.Output.Body([]byte(fileURL))
}

func main() {
	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		return
	}

	// 配置日志
	err := logs.SetLogger(logs.AdapterFile, `{"filename":"logs/app.log","level":7,"maxlines":10000,"maxsize":0,"daily":true,"maxdays":10}`)
	if err != nil {
		fmt.Printf("Failed to set logger: %v\n", err)
		return
	}
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.Async()

	// 启动机器人
	bot.Start(conf.TelegramBotToken)

	// 注册路由
	web.Router("/:file_id", &FileController{})

	// 启动 Beego 服务器，监听 8001 端口
	web.Run(":8001")
}
