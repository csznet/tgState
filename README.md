# tgstate

使用Telegram作为存储的图床，最大支持20MB

搭配CloudFlare：https://www.csz.net/proj/tgstate/

如果需要更多参数支持，建议使用二进制方式运行

如有疑惑，可以咨询TG @tgstate123

自夸：
 - 原图上传
 - 支持Vercel一键搭建
 - 支持粘贴上传
 - 支持一键复制bbcode markdown html代码
 - 储存在自己的频道里，最大支持20MB
 - 提供API接口，方便二次开发

先决条件：
 - 创建机器人获取Bot Token

关于channel，可以设置为频道，格式为@xxxx，也可以设置为自己的ID
设置频道需要创建频道，将机器人拉入作为管理员（公开频道，私有频道没有测试）
设置为telegram的user id需要先给机器人发一条信息，telegram user id查看（@getmyid_bot）

建议套上CloudFlare 设置 ```/d/*```完全缓存，并开启always online  

测试地址：https://tgtu.ikun123.com/  (搭建在vercel)  
测试图片：  
![tgState](https://tgtu.ikun123.com/img/364.jpg)  


Vercel安装
====

[点我传送至Vercel配置页面](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fcsznet%2FtgState&env=token&env=channel&project-name=tgState&repository-name=tgState)  
token填写你的bot token，channel可以为频道(@xxxx)，也可以为你的telegram id(@getmyid_bot获取)


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

请提前将aaa和bbb替换为你的bot token和频道地址or个人id


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

其中的xxxx为bot token @xxxx为频道地址or个人id

如果需要自定义端口，可以带上-port参数，如
```
-port 8888
```
如果不需要首页，只需要API和图片展示页面，则带上-index参数，如
```
 ./tgState -token xxxx -channel @xxxx -port 8888 -index
 ```  