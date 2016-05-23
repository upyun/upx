> upx is a tool for managing files in UPYUN. Mac, Linux, Windows supported

[![Build Status](https://travis-ci.org/polym/upx.svg?branch=master)](https://travis-ci.org/polym/upx)

## 基本功能

- [x] 支持基本文件系统操作命令，如 `mkdir`, `cd`, `ls`, `rm`, `pwd`
- [x] 支持上传文件或目录到又拍云存储
- [x] 支持从又拍云存储下载文件或目录到本地
- [x] 支持增量同步文件到又拍云存储
- [x] 支持删除又拍云存储中的文件或目录，并且支持通配符 `*`
- [x] 支持多用户，多操作系统

## 安装

### 安装包下载地址

- [Linux 64位](http://collection.b0.upaiyun.com/softwares/upx/upx-linux-amd64-v0.1.3)
- [Linux 32位](http://collection.b0.upaiyun.com/softwares/upx/upx-linux-i386-v0.1.3)
- [Windows 64位](http://collection.b0.upaiyun.com/softwares/upx/upx-windows-amd64-v0.1.3.exe)
- [Windows 32位](http://collection.b0.upaiyun.com/softwares/upx/upx-windows-i386-v0.1.3.exe)
- [Mac 64位](http://collection.b0.upaiyun.com/softwares/upx/upx-darwin-amd64-v0.1.3)

### 源码编译

> 需要安装 [Golang 编译环境](https://golang.org/dl/)

```
$ git clone https://github.com/polym/upx.git
$ cd upx && make
```
or

```
$ go get github.com/polym/upx
```

## 使用

> 所有命令都支持 `-h` 查看使用方法

|  命令  | 说明 |
| ------ | ---- |
| login  | 登录又拍云存储 |
| logout | 退出帐号 |
| mkdir  | 创建目录 |
| pwd    | 显示当前所在目录 |
| ls     | 显示当前目录下文件和目录信息 |
| info   | 显示服务名、用户名等信息 |
| cd     | 改变工作目录（进入一个目录）|
| get    | 下载一个文件或目录 |
| put    | 上传一个文件或目录 |
| sync   | 目录增量同步，类似 rsync |
| rm     | 删除目录或文件 |


| global options | 说明 |
| -------------- | ---- |
| -h             | 显示帮助信息 |
| -v             | 显示 UPX 版本信息 |


### 列目录 `ls`

> 默认按文件修改时间先后顺序输出

| options | 说明 |
| ------- | ---- |
| -d      | 仅显示目录 |
| -r      | 文件修改时间倒序输出 |
| -c v    | 仅显示前 v 个文件或目录  |

### 删除 `rm`

> 默认不会删除目录，支持通配符 `*`

| options | 说明 |
| ------- | ---- |
| -d      | 仅删除目录 |
| -a      | 删除目录跟文件 |
| --async | 异步删除，目录可能需要二次删除 |

### 增量同步 `sync`

> sync 本地路径 存储路径

| options | 说明 |
| ------- | ---- |
| -v      | 是否显示详细信息 |
| -w      | 制定并发数，默认为 10 |
