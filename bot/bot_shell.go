package bot

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"tg-file-share/conf"

	"github.com/beego/beego/v2/core/logs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	running  bool
	stopChan chan struct{}
	mu       sync.Mutex
	updates  tgbotapi.UpdatesChannel
}

// New 创建新的机器人实例
func New(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	// 添加调试模式
	api.Debug = true

	bot := &Bot{
		api: api,
	}

	return bot, nil
}

// Start 启动机器人
func Start(token string) error {
	bot, err := New(token)
	if err != nil {
		logs.Error("创建机器人失败: %s", err)
		return err
	}

	go bot.Start()
	return nil
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	b.updates = b.api.GetUpdatesChan(u)

	for {
		select {
		case <-b.stopChan:
			log.Printf("消息转发Bot 已停止...")
			return nil
		case update, ok := <-b.updates:
			if !ok {
				return nil
			}
			if update.Message == nil {
				continue
			}

			log.Printf("[%s] 收到消息: MessageID: %s (from: %s, chat_id: %d)",
				update.Message.MessageID,
				update.Message.Text,
				update.Message.From.UserName,
				update.Message.Chat.ID)

			if update.Message.IsCommand() {
				log.Printf("[%s] 命令消息: %s", update.Message.Command())
				continue
			}

			b.handleCommand(update.Message)
		}
	}
}

// handleCommand 处理命令消息
func (b *Bot) handleCommand(message *tgbotapi.Message) {
	// 检查是否有任何内容需要处理
	hasContent := message.Text != "" ||
		message.Sticker != nil ||
		message.Animation != nil ||
		message.Video != nil ||
		message.Location != nil ||
		message.Poll != nil ||
		message.Document != nil ||
		message.Photo != nil ||
		message.Voice != nil

	if !hasContent {
		return
	}
	// targetChatID 改成谁跟机器人说话，就回复谁
	// targetChatID = message.Chat.ID

	// 处理文本消息
	if message.Text != "" {
		// forwardText += message.Text
		// msg := tgbotapi.NewMessage(targetChatID, forwardText)
		// msg.ParseMode = tgbotapi.ModeMarkdownV2 // 设置 Markdown V2 解析模式
		// b.sendWithLog(msg, "text message")
		return
	}

	// 处理文档（包括 GIF）
	if message.Document != nil {
		// 直接回复用户：tg.iamxmm.xyz/message.Document.FileID
		// 如何直接在原文的下面回复用户
		forwardText := conf.FileServerURL + "/" + message.Document.FileID
		b.ReplyToMessage(message, forwardText)
	}

	// 处理图片
	if message.Photo != nil && len(message.Photo) > 0 {
		photo := message.Photo[len(message.Photo)-1]
		forwardText := conf.FileServerURL + "/" + photo.FileID
		b.ReplyToMessage(message, forwardText)

		// photo := message.Photo[len(message.Photo)-1]
		// photoMsg := tgbotapi.NewPhoto(targetChatID, tgbotapi.FileID(photo.FileID))
		// photoMsg.Caption = senderInfo
		// b.sendWithLog(photoMsg, "photo")
	}

	// 处理视频
	if message.Video != nil {
		// videoMsg := tgbotapi.NewVideo(targetChatID, tgbotapi.FileID(message.Video.FileID))
		// videoMsg.Caption = senderInfo
		// b.sendWithLog(videoMsg, "video")
		forwardText := conf.FileServerURL + "/" + message.Video.FileID
		b.ReplyToMessage(message, forwardText)
	}
}

// sendWithLog 统一处理消息发送和错误日志
func (b *Bot) sendWithLog(msg tgbotapi.Chattable, msgType string) {
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Failed to forward %s: %v", msgType, err)
	}
	log.Printf("消息【%s】发送成功", msgType)
}

// 检查文件是否是 GIF
func isGif(fileName string) bool {
	if fileName == "" {
		return false
	}
	return strings.ToLower(filepath.Ext(fileName)) == ".gif"
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// 在原消息下方回复用户
func (b *Bot) ReplyToMessage(message *tgbotapi.Message, text string) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID // 设置回复的消息ID，这样回复会显示在原消息下方

	_, err := b.api.Send(msg)
	if err != nil {
		logs.Error("Failed to send reply message: %v", err)
		return err
	}

	return nil
}
