package utils

import (
	"encoding/json"
	"io"
	"log"

	"csz.net/tgstate/conf"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TgFileData(fileName string, fileData io.Reader) tgbotapi.FileReader {
	return tgbotapi.FileReader{
		Name:   fileName,
		Reader: fileData,
	}
}

func UpDocument(fileData tgbotapi.FileReader) string {
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	// Upload the file to Telegram
	params := tgbotapi.Params{
		"chat_id": conf.ChannelName, // Replace with the chat ID where you want to send the file
	}
	files := []tgbotapi.RequestFile{
		{
			Name: "document",
			Data: fileData,
		},
	}

	response, err := bot.UploadFiles("sendDocument", params, files)
	if err != nil {
		log.Panic(err)
	}
	var msg tgbotapi.Message
	err = json.Unmarshal([]byte(response.Result), &msg)
	var resp string
	if msg.Document != nil {
		resp = msg.Document.FileID
	} else if msg.Audio != nil {
		resp = msg.Audio.FileID
	} else if msg.Video != nil {
		resp = msg.Video.FileID
	}
	return resp
}

func GetDownloadUrl(fileID string) string {
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Panic(err)
	}

	// 使用 getFile 方法获取文件信息
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		log.Panic(err)
	}

	// 获取文件下载链接
	fileURL := file.Link(conf.BotToken)
	// log.Printf("File Download URL: %s", fileURL)
	return fileURL
}
func BotDo() {
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Println(err)
		return
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updatesChan := bot.GetUpdatesChan(u)
	for update := range updatesChan {
		var msg *tgbotapi.Message
		if update.Message != nil {
			msg = update.Message
		}
		if update.ChannelPost != nil {
			msg = update.ChannelPost
		}
		if msg != nil && msg.Text == "get" && msg.ReplyToMessage != nil {
			var fileID string
			if msg.ReplyToMessage.Document != nil && msg.ReplyToMessage.Document.FileID != "" {
				fileID = msg.ReplyToMessage.Document.FileID
			}
			if msg.ReplyToMessage.Video != nil && msg.ReplyToMessage.Video.FileID != "" {
				fileID = msg.ReplyToMessage.Video.FileID
			}
			if fileID != "" {
				newMsg := tgbotapi.NewMessage(msg.Chat.ID, fileID)
				newMsg.ReplyToMessageID = msg.MessageID
				bot.Send(newMsg)
			}
		}
	}
}
