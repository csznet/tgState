tgState
==

[中文](https://github.com/csznet/tgState/blob/main/README.md)


A file external link system that uses Telegram as storage.

It can be used as a Telegram image hosting service or as a Telegram cloud storage.

The default operating mode is image hosting mode, which only allows the upload of .jpg, .png, and .jpeg files with a maximum size limit of 20MB. In contrast, the cloud storage mode has no restrictions on file extensions or file sizes.

If you have any questions or concerns, feel free to reach out to us on Telegram at @tgstate123.

**Features**
 - No file size limit (optional)
 - Support for access password restrictions
 - Provides API
 - Supports one-click deployment on Vercel

**Demo**

https://tgtu.ikun123.com/  
Hosted on Vercel, large files may fail to upload 

Test image：

![tgState](https://tgstate.vercel.app/d/BQACAgUAAx0EcyK3ugACByxlOR-Nfl4esavoO4zdaYIP_k1KYQACDAsAAkf4yFVpf_awaEkS8jAE)  

**Preparation Instructions

**
Before deployment, you need to prepare a Telegram Bot (apply at @botfather).
If you need to store files in a channel, you need to add the Bot to the channel as an administrator, make the channel public, and customize the channel link.

Vercel Deployment
====

 [Click here to go to the Vercel configuration page](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fcsznet%2FtgState&env=token&env=channel&env=pass&env=mode&project-name=tgState&repository-name=tgState)  

 1. Fill in your bot `token`
 2. `channel`` can be a channel (@xxxx) or your Telegram ID (use @getmyid_bot to obtain)  
 3. Fill in the access password for `pass`, if not needed, simply enter `none`
 4. Fill in `mode` as `pan` to run in cloud storage mode, or enter anything for image hosting mode   

Docker Deployment
====

Pull the image:
```
docker pull csznet/tgstate:latest
```

Start the container:
```
docker run -d -p 8088:8088 --name tgstate -e TOKEN=aaa -e CHANNEL=@bbb csznet/tgstate:latest
```

Replace aaa and bbb with your bot token and channel address or personal ID.  

If you need to run in cloud storage mode:

```
docker run -d -p 8088:8088 --name tgstate -e TOKEN=aaa -e CHANNEL=@bbb -e MODE=pan csznet/tgstate:latest
```


 Binary Deployment
====
Download the binary file for Linux amd64:
 ```
 wget https://github.com/csznet/tgState/releases/latest/download/tgState.zip
 ```

Unzip the file:


 
 ```
 unzip tgState.zip && rm tgState.zip
 ```
Usage
----

```
 ./tgState -token xxxx -channel @xxxx
```

Replace `xxxx` with your bot token and `@xxxx` with the channel address or personal ID (for personal IDs, only use numbers without @).

To customize the port, use the `-port` parameter, for example:
```
-port 8888
```
If you don't need the homepage and only want the API and image display page, use the `-index` parameter, like this:
```
./tgState -token xxxx -channel @xxxx -port 8888 -index
```  
To enable password protection, use the `-pass` parameter, for example, to set the password as `csznet`:
```
./tgState -token xxxx -channel @xxxx -port 8888 -pass csznet
```

For cloud storage mode, use the `-mode pan` parameter, like this:

```
./tgState -token xxxx -channel @xxxx -port 8888 -mode pan
```

About the API   
====

Send file data directly as binary to the `/api` path.