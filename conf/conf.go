package conf

var BotToken string
var ChannelName string
var Pass string
var Mode string
var BaseUrl string
var AllowedExts string
var ProxyUrl string

type UploadResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	ImgUrl   string `json:"url"`
	ProxyUrl string `json:"proxyUrl"`
}

const FileRoute = "/d/"
