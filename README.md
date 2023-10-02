# tgstate

使用Telegram作为存储的图床，最大支持20MB

言简意赅的教程：https://www.csz.net/proj/tgstate/

自夸：
 - 支持粘贴上传
 - 支持一键复制bbcode markdown html代码
 - 储存在自己的频道里，最大支持20MB
 - 提供API接口，方便二次开发

先决条件：
 - 创建机器人获取Bot Token
 - 创建频道，将机器人拉入作为管理员（公开频道，私有频道没有测试）

建议套上CloudFlare 设置 ```/img/*``` 和```/d/*```完全缓存，并开启always online  

测试地址：https://tgtu.ikun123.com/  
测试图片：![tgState](https://tgtu.ikun123.com/img/364.jpg)


Docker安装
====

pull镜像
```
docker pull csznet/tgstate:latest
```

启动
```
docker run -d -p 8088:8088 --name tgstate -e TOKEN=aaa -e CHANNEL=@bbb csznet/tgstate:latest
```

请提前将aaa和bbb替换为你的bot token和频道地址


 二进制安装
====
 下载Linux amd64环境的二进制文件
 ```
 wget https://github.com/csznet/tgState/releases/latest/download/tgState.zip
 ```
 解压
 ```
 unzip tgState.zip && rm tgState.zip
 ```
 使用方法
----

```
 ./tgState -token xxxx -channel @xxxx
```

其中的xxxx为bot token @xxxx为频道地址

如果需要自定义端口，可以带上-port参数，如
```
-port 8888
```
如果不需要首页，只需要API和图片展示页面，则带上-index参数，如
```
 ./tgState -token xxxx -channel @xxxx -port 8888 -index
 ```
