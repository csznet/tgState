package api

import (
	"net/http"
	"os"
	"strings"

	"csz.net/tgstate/conf"
	"csz.net/tgstate/control"
)

func Vercel(w http.ResponseWriter, r *http.Request) {
	conf.BotToken = os.Getenv("token")
	conf.ChannelName = os.Getenv("channel")
	// 获取请求路径
	path := r.URL.Path
	// 如果请求路径以 "/img/" 开头
	if strings.HasPrefix(path, "/img/") {
		// 调用 control 包中的 Img 处理函数
		control.Img(w, r)
		return // 结束处理，确保不执行默认处理
	}
	if strings.HasPrefix(path, "/d/") {
		// 调用 control 包中的 Img 处理函数
		control.D(w, r)
		return // 结束处理，确保不执行默认处理
	}
	switch path {
	case "/api":
		// 调用 control 包中的 UploadImageAPI 处理函数
		control.UploadImageAPI(w, r)
	default:
		control.Index(w, r)
	}
}
