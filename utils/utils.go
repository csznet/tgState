package utils

import (
	"encoding/json"
	"log"

	"csz.net/tgstate/conf"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TgFileData(fileName string, fileData []byte) tgbotapi.FileBytes {
	return tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: fileData,
	}
}

func UpDocument(fileData tgbotapi.FileBytes) string {
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
	var msg conf.Message
	err = json.Unmarshal([]byte(response.Result), &msg)
	var resp string
	if msg.Document.FileID != "" {
		resp = msg.Document.FileID
	} else if msg.Audio.FileID != "" {
		resp = msg.Audio.FileID
	} else if msg.Video.FileID != "" {
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
