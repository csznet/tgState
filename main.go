package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"csz.net/tgstate/conf"
	"csz.net/tgstate/control"
	"csz.net/tgstate/utils"
)

var webPort string
var index = true

func main() {
	//判断是否设置参数
	if conf.BotToken == "" || conf.ChannelName == "" {
		fmt.Println("请先设置Bot Token和对象")
		return
	}
	go utils.BotDo()
	web()
}

func web() {
	if conf.Pass != "" && conf.Pass != "none" {
		http.HandleFunc("/pwd", control.Pwd)
	}
	http.HandleFunc("/d/", control.D)
	http.HandleFunc("/api", control.Middleware(control.UploadImageAPI))
	if index {
		http.HandleFunc("/", control.Middleware(control.Index))
	}
	listener, err := net.Listen("tcp", ":"+webPort)
	if err != nil {
		fmt.Printf("端口 %s 已被占用\n", webPort)
		return
	}
	defer listener.Close()
	fmt.Printf("启动Web服务器，监听端口 %s\n", webPort)
	err = http.Serve(listener, nil)
	if err != nil {
		fmt.Println(err)
	}
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
