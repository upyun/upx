> upx is a tool for managing files in UPYUN. Mac, Linux, Windows supported

![Test](https://github.com/upyun/upx/workflows/Test/badge.svg)
![Build](https://github.com/upyun/upx/workflows/Build/badge.svg)
![Lint](https://github.com/upyun/upx/workflows/Lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/upyun/upx)](https://goreportcard.com/report/github.com/upyun/upx)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/upyun/upx?label=latest%20release)

## 基本功能

- [x] 支持基本文件系统操作命令，如 `mkdir`, `cd`, `ls`, `rm`, `pwd`
- [x] 支持上传文件或目录到又拍云存储
- [x] 支持从又拍云存储下载文件或目录到本地
- [x] 支持增量同步文件到又拍云存储
- [x] 支持删除又拍云存储中的文件或目录，并且支持通配符 `*`
- [x] 支持多用户，多操作系统
- [x] 支持基于时间列目录以及删除文件
- [x] 支持 `tree` 获取目录结构
- [x] 支持提交异步处理任务
- [x] 更加准确简洁的进度条
- [x] 使用 UPYUN GoSDK v3
- [x] 同步目录支持 --delete
- [x] 支持 CDN 缓存刷新

## 安装

### 可执行程序二进制下载地址

- [Windows x86_64](http://collection.b0.upaiyun.com/softwares/upx/upx_0.3.6_Windows_x86_64.zip)
- [Windows i386](http://collection.b0.upaiyun.com/softwares/upx/upx_0.3.6_Windows_i386.zip)
- [Mac x86_64](http://collection.b0.upaiyun.com/softwares/upx/upx_0.3.6_Darwin_x86_64.tar.gz)
- [Linux x86_64](http://collection.b0.upaiyun.com/softwares/upx/upx_0.3.6_linux_x86_64.tar.gz)
- [Linux i386](http://collection.b0.upaiyun.com/softwares/upx/upx_0.3.6_linux_i386.tar.gz)
- [Linux arm64](http://collection.b0.upaiyun.com/softwares/upx/upx_0.3.6_linux_arm64.tar.gz)
- [Linux armv6](http://collection.b0.upaiyun.com/softwares/upx/upx_0.3.6_linux_armv6.tar.gz)
- [Linux armv7](http://collection.b0.upaiyun.com/softwares/upx/upx_0.3.6_linux_armv7.tar.gz)

### 源码编译

> 需要安装 [Golang 编译环境](https://golang.org/dl/)

```
$ git clone https://github.com/upyun/upx.git
$ cd upx && make
```
or

```
$ GO111MODULE=on go get -u github.com/upyun/upx@v0.3.6
```

### Windows

```
PS> scoop bucket add carrot https://github.com/upyun/carrot.git
Install upx from github or upyun cdn:
PS> scoop install upx-github
PS> scoop install upx-upcdn
```

### Docker

```bash
docker build -t upx .
docker run --rm upx upx -v
```

## 使用

> 所有命令都支持 `-h` 查看使用方法

|    命令  | 说明 |
| -------- | ---- |
| login    | 登录又拍云存储 |
| logout   | 退出帐号 |
| sessions | 查看所有的会话 |
| switch   | 切换会话 |
| info     | 显示服务名、用户名等信息 |
| cd       | 改变工作目录（进入一个目录）|
| pwd      | 显示当前所在目录 |
| mkdir    | 创建目录 |
| ls       | 显示当前目录下文件和目录信息 |
| tree     | 显示目录结构 |
| get      | 下载一个文件或目录 |
| put      | 上传一个文件或目录 |
| rm       | 删除目录或文件 |
| sync     | 目录增量同步，类似 rsync |
| auth     | 生成包含空间名操作员密码信息的 auth 字符串 |
| post     | 提交异步处理任务 |
| purge    | 提交 CDN 缓存刷新任务 |


| global options | 说明 |
| -------------- | ---- |
| --quiet, -q    | 不显示信息 |
| --auth value   | auth 字符串 |
| --help, -h     | 显示帮助信息 |
| --version, -v  | 显示版本号 |


### 列目录 `ls`

> 默认按文件修改时间先后顺序输出

|  options  | 说明 |
| --------- | ---- |
| -d        | 仅显示目录 |
| -r        | 文件修改时间倒序输出 |
| --color   | 根据文件类型输出不同的颜色 |
| -c v      | 仅显示前 v 个文件或目录, 默认全部显示  |
| --mtime v | 参考 Linux `find` |

### 删除 `rm`

> 默认不会删除目录，支持通配符 `*`

|  options  | 说明 |
| --------- | ---- |
| -d        | 仅删除目录 |
| -a        | 删除目录跟文件 |
| --async   | 异步删除，目录可能需要二次删除 |
| --mtime v | 参考 Linux `find` |

### 增量同步 `sync`

> sync 本地路径 存储路径

| options  | 说明 |
| -------- | ---- |
| -w       | 指定并发数，默认为 5 |
| --delete | 删除上一次同步后本地删除的文件 |

### 生成 auth 串 `auth`

> auth 空间名 操作员 密码

当命令中包含 `--auth` 参数时，会忽略已登陆的信息。

### 提交异步任务 `post`

|     options    | 说明 |
| -------------- | ---- |
| --app value    | app 名称 |
| --notify value | 回调地址 |
| --task value   | 任务文件名 |

### 缓存刷新 `purge`

> purge url --list urls

|     options    | 说明 |
| -------------- | ---- |
| --list value   | 批量刷新文件名 |

## TODO

- [ ] put 支持断点续传
- [ ] upx 支持指定 API 地址
