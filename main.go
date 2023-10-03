package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"csz.net/tgstate/conf"
	"csz.net/tgstate/control"
)

var webPort string
var index = true

func main() {
	//判断是否设置参数
	if conf.BotToken == "" || conf.ChannelName == "" {
		fmt.Println("请先设置Bot Token和对象")
		return
	}
	web()
}

func web() {
	http.HandleFunc("/img/", control.Img)
	http.HandleFunc("/d/", control.D)
	http.HandleFunc("/api", control.UploadImageAPI)
	if index {
		http.HandleFunc("/", control.Index)
	}
	http.ListenAndServe(":"+webPort, nil)
}

func init() {
	flag.StringVar(&webPort, "port", "8088", "Web Port")
	flag.StringVar(&conf.BotToken, "token", "", "Bot Token")
	flag.StringVar(&conf.ChannelName, "channel", "", "Channel Name")
	indexPtr := flag.Bool("index", false, "Show Index")
	flag.Parse()
	if *indexPtr {
		index = false
	}
	if conf.BotToken == "" {
		conf.BotToken = os.Getenv("TOKEN")
	}
	if conf.ChannelName == "" {
		conf.ChannelName = os.Getenv("CHANNEL")
	}
}
