@echo off
REM 设置目标操作系统为 Linux，架构为 64 位
set GOOS=freebsd
set GOARCH=amd64

REM 编译 Go 程序并输出为指定的二进制文件名
go build -o tgState main.go

REM 提示编译完成
echo 编译完成，生成了 Linux 版本的二进制文件 tgState
pause
