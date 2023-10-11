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
	//输出版本号
	fmt.Printf("版本号: %s\n", "v1.2.2")
	//判断是否设置参数
	if conf.BotToken == "" || conf.ChannelName == "" {
		fmt.Println("请先设置Bot Token和对象")
		return
	}
	web()
}

func web() {
	http.HandleFunc("/pwd", control.Pwd)
	http.HandleFunc("/d/", control.D)
	http.HandleFunc("/api", control.Middleware(control.UploadImageAPI))
	if index {
		http.HandleFunc("/", control.Middleware(control.Index))
	}
	http.ListenAndServe(":"+webPort, nil)
}

func init() {
	flag.StringVar(&webPort, "port", "8088", "Web Port")
	flag.StringVar(&conf.BotToken, "token", "", "Bot Token")
	flag.StringVar(&conf.ChannelName, "channel", "", "Channel Name")
	flag.StringVar(&conf.Pass, "pass", "", "Visit Password")
	flag.StringVar(&conf.Mode, "mode", "", "Run mode")
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
	if conf.Mode == "" {
		conf.Mode = os.Getenv("MODE")
	}
}
