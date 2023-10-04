package conf

var BotToken string
var ChannelName string
var Pass string

type Thumbnail struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

type Document struct {
	FileName     string    `json:"file_name"`
	MimeType     string    `json:"mime_type"`
	Thumbnail    Thumbnail `json:"thumbnail"`
	Thumb        Thumbnail `json:"thumb"`
	FileID       string    `json:"file_id"`
	FileUniqueID string    `json:"file_unique_id"`
	FileSize     int       `json:"file_size"`
}

type SenderChat struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Username string `json:"username"`
	Type     string `json:"type"`
}

type Chat struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Username string `json:"username"`
	Type     string `json:"type"`
}

type Message struct {
	MessageID  int64      `json:"message_id"`
	SenderChat SenderChat `json:"sender_chat"`
	Chat       Chat       `json:"chat"`
	Date       int64      `json:"date"`
	Document   Document   `json:"document"`
}
type UploadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
