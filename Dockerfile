# 使用官方的 Ubuntu 基础镜像
FROM ubuntu:latest

# 安装 ca-certificates 包，用于更新根证书
RUN apt-get update && apt-get install -y ca-certificates

# 将编译好的 server 和 client 二进制文件复制到容器中
COPY tgState /app/tgState

# 设置工作目录
WORKDIR /app

# 设置暴露的端口
EXPOSE 8088

# 设置容器启动时要执行的命令
CMD  [ "/app/tgState" ]