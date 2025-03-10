package conf

import "github.com/beego/beego/v2/server/web"

// Telegram API 相关信息
var (
	TelegramBotToken = ""
	FileServerURL    = ""
	TelegramAPI      = "https://api.telegram.org"
)

func init() {
	TelegramBotToken = web.AppConfig.DefaultString("telegram_bot_token", "")
	if TelegramBotToken == "" {
		panic("telegram_bot_token 不能为空")
	}
	FileServerURL = web.AppConfig.DefaultString("file_server_url", "")
	if FileServerURL == "" {
		panic("file_server_url 不能为空")
	}
}
