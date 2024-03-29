# tgState
==

[中文](https://github.com/csznet/tgState/blob/main/README.md)

A file external link system using Telegram as storage.

No restrictions on file size and format.

Can be used as a Telegram image hosting service or a Telegram cloud drive.

Supports web and Telegram direct file uploads.

Use with CloudFlare: https://github.com/csznet/tgState/blob/main/CloudFlare.md

For any questions, consult TG @tgstate123

# Demo

https://tgstate.vercel.app / https://tgstate.ikun123.com/

Hosted on Vercel, resource limitations - files larger than 5MB are not supported.

Demo image:

![tgState](https://tgstate.vercel.app/d/BQACAgUAAx0EcyK3ugACByxlOR-Nfl4esavoO4zdaYIP_k1KYQACDAsAAkf4yFVpf_awaEkS8jAE)

# Parameter Description

Mandatory parameters:

- target
- token

Optional parameters:

- pass
- mode
- url
- port

## target

The target can be a channel, group, or individual.

When the target is a channel, the bot needs to be added to the channel as an administrator, make the channel public, and customize the channel link. The target value should be filled with the link, such as @xxxx.

When the target is a group, the bot needs to be added to the group, make the group public, and customize the group link. The target value should be filled with the link, such as @xxxx.

When the target is an individual, it is the Telegram ID (obtained from @getmyid_bot).

## token

Fill in your bot token.

## pass

Fill in the access password. If not needed, fill in ```none``` directly.

## mode

- ```p``` represents running in cloud drive mode, with no restriction on uploaded suffixes.
- ```m``` On top of the p mode, web upload is disabled, and upload can be done via private chat (if the target is an individual, only specified users can upload via private chat).

## url

The pre-domain address that bot obtains FileID is automatically filled in.

## port

Customize the running port.

# Management

## Get FIleID

Replying with ```get``` to the file reference in the bot's chat can get the FileID. Access the resource by combining the built address and the obtained path.

If the url parameter is configured, the complete address will be returned directly.

![image](https://github.com/csznet/tgState/assets/127601663/5b1fd6c0-652c-41de-bb63-e2f20b257022)

# Deployment

## Binary

Download for Linux amd64

```
wget https://github.com/csznet/tgState/releases/latest/download/tgState.zip && unzip tgState.zip && rm tgState.zip
```

Download for Linux arm64

```
wget https://github.com/csznet/tgState/releases/latest/download/tgState_arm64.zip && unzip tgState_arm64.zip && rm tgState_arm64.zip
```

**Usage**

```./tgState parameters```

**Example**

```./tgState -token xxxx -target @xxxx```

**Run in the background**

```nohup ./tgState parameters &```

## Docker

Pull the image

```docker pull csznet/tgstate:latest```


Start

```
docker run -d -p 8088:8088 --name tgstate parameters --net=host csznet/tgstate:latest
```

Where docker parameters need to be set as environment variables.

**Example**

```
docker run -d -p 8088:8088 --name tgstate -e token=aaa -e target=@bbb --net=host csznet/tgstate:latest
```

## Vercel

Does not support files larger than 5MB and does not support Telegram in getting file paths.

[Click here to go to the Vercel configuration page](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fcsznet%2FtgState&env=token&env=target&env=pass&env=mode&project-name=tgState&repository-name=tgState)

# API Description

POST method to the path ```/api```

Form transmission, field name is image, content is binary data.
