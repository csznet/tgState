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
var OptApi = true

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
	http.HandleFunc(conf.FileRoute, control.D)
	if OptApi {
		if conf.Pass != "" && conf.Pass != "none" {
			http.HandleFunc("/pwd", control.Pwd)
		}
		http.HandleFunc("/api", control.Middleware(control.UploadImageAPI))
		http.HandleFunc("/", control.Middleware(control.Index))
	}

	if listener, err := net.Listen("tcp", ":"+webPort); err != nil {
		fmt.Printf("端口 %s 已被占用\n", webPort)
	} else {
		defer listener.Close()
		fmt.Printf("启动Web服务器，监听端口 %s\n", webPort)
		if err := http.Serve(listener, nil); err != nil {
			fmt.Println(err)
		}
	}
}

func init() {
	flag.StringVar(&webPort, "port", "8088", "Web Port")
	flag.StringVar(&conf.BotToken, "token", os.Getenv("token"), "Bot Token")
	flag.StringVar(&conf.ChannelName, "target", os.Getenv("target"), "Channel Name or ID")
	flag.StringVar(&conf.Pass, "pass", os.Getenv("pass"), "Visit Password")
	flag.StringVar(&conf.Mode, "mode", os.Getenv("mode"), "Run mode")
	flag.StringVar(&conf.BaseUrl, "url", os.Getenv("url"), "Base Url")
	flag.Parse()
	if conf.Mode == "m" {
		OptApi = false
	}
	if conf.Mode != "p" && conf.Mode != "m" {
		conf.Mode = "p"
	}
}
