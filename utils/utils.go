package utils

import (
	"encoding/json"
	"io"
	"log"
	"strconv"
	"strings"

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
		log.Println(err)
		return ""
	}
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
	json.Unmarshal([]byte(response.Result), &msg)
	var resp string
	switch {
	case msg.Document != nil:
		resp = msg.Document.FileID
	case msg.Audio != nil:
		resp = msg.Audio.FileID
	case msg.Video != nil:
		resp = msg.Video.FileID
	case msg.Sticker != nil:
		resp = msg.Sticker.FileID
	}
	return resp
}

func GetDownloadUrl(fileID string) (string, bool) {
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Panic(err)
	}
	// 使用 getFile 方法获取文件信息
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		log.Println("获取文件失败【" + fileID + "】")
		log.Println(err)
		return "", false
	}
	log.Println("获取文件成功【" + fileID + "】")
	// 获取文件下载链接
	fileURL := file.Link(conf.BotToken)
	return fileURL, true
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
			switch {
			case msg.ReplyToMessage.Document != nil && msg.ReplyToMessage.Document.FileID != "":
				fileID = msg.ReplyToMessage.Document.FileID
			case msg.ReplyToMessage.Video != nil && msg.ReplyToMessage.Video.FileID != "":
				fileID = msg.ReplyToMessage.Video.FileID
			case msg.ReplyToMessage.Sticker != nil && msg.ReplyToMessage.Sticker.FileID != "":
				fileID = msg.ReplyToMessage.Sticker.FileID
			}
			if fileID != "" {
				newMsg := tgbotapi.NewMessage(msg.Chat.ID, strings.TrimSuffix(conf.BaseUrl, "/")+"/d/"+fileID)
				newMsg.ReplyToMessageID = msg.MessageID
				if !strings.HasPrefix(conf.ChannelName, "@") {
					if man, err := strconv.Atoi(conf.ChannelName); err == nil && int(msg.Chat.ID) == man {
						bot.Send(newMsg)
					}
				} else {
					bot.Send(newMsg)
				}
			}
		}
	}
}
