tgState
==

[English](https://github.com/csznet/tgState/blob/main/README_en.md)


一款以Telegram作为储存的文件外链系统

可以作为telegram图床，也可以作为telegram网盘使用。

默认运行模式为图床模式，只允许`.jpg .png .jpeg`文件上传且限制不超过20MB，网盘模式为不限制后缀和大小  

如有疑惑，可以咨询TG @tgstate123  

**特性**
 - 不限制上传文件大小（可选
 - 支持访问密码限制
 - 提供API
 - 支持Vercel一键搭建

**Demo**

https://tgtu.ikun123.com/  
搭建在Vercel，大文件可能会上传失败  

测试图片：

![tgState](https://tgtu.ikun123.com/d/1310.png)  

**准备说明**
部署前需要准备一个Telegram Bot(@botfather处申请)  
如果是需要存储在频道，则需要将Bot拉入频道作为管理员，公开频道并自定义频道Link  

Vercel部署
====

 [点我传送至Vercel配置页面](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fcsznet%2FtgState&env=token&env=channel&env=pass&env=mode&project-name=tgState&repository-name=tgState)  

 1. ```token```填写你的bot token  
 2. ```channel```可以为频道(@xxxx)，也可以为你的telegram id(@getmyid_bot获取)  
 3. ```pass```填写访问密码，如不需要，直接填写```none```即可
 4. ```mode```填写```pan```，代表以网盘模式运行,只需要以图床模式运行的话就随便填    

 Docker部署
====

pull镜像
```
docker pull csznet/tgstate:latest
```

启动
```
docker run -d -p 8088:8088 --name tgstate -e TOKEN=aaa -e CHANNEL=@bbb csznet/tgstate:latest
```

请提前将```aaa```和```bbb```替换为你的bot token和频道地址or个人id  

如果需要以网盘模式启动  

```
docker run -d -p 8088:8088 --name tgstate -e TOKEN=aaa -e CHANNEL=@bbb -e MODE=pan csznet/tgstate:latest
```


 二进制部署
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

其中的```xxxx```为bot token ```@xxxx```为频道地址or个人id(个人ID只需要数字不需要@)

如果需要自定义端口，可以带上-port参数，如
```
-port 8888
```
如果不需要首页，只需要API和图片展示页面，则带上-index参数，如
```
./tgState -token xxxx -channel @xxxx -port 8888 -index
```  
如果需要限制密码访问，只需要带上-pass参数即可，如设置密码为csznet：  
```
./tgState -token xxxx -channel @xxxx -port 8888 -pass csznet
```

如果需要网盘模式运行，请带上-mode pan，如  

```
./tgState -token xxxx -channel @xxxx -port 8888 -mode pan
```

关于API  
====

直接将文件数据以二进制的方式发送给```/api```路径