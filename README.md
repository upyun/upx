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

- [Windows x86_64](https://collection.b0.upaiyun.com/softwares/upx/upx_0.3.7_Windows_x86_64.zip)
- [Windows i386](https://collection.b0.upaiyun.com/softwares/upx/upx_0.3.7_Windows_i386.zip)
- [Mac x86_64](https://collection.b0.upaiyun.com/softwares/upx/upx_0.3.7_Darwin_x86_64.tar.gz)
- [Linux x86_64](https://collection.b0.upaiyun.com/softwares/upx/upx_0.3.7_linux_x86_64.tar.gz)
- [Linux i386](https://collection.b0.upaiyun.com/softwares/upx/upx_0.3.7_linux_i386.tar.gz)
- [Linux arm64](https://collection.b0.upaiyun.com/softwares/upx/upx_0.3.7_linux_arm64.tar.gz)
- [Linux armv6](https://collection.b0.upaiyun.com/softwares/upx/upx_0.3.7_linux_armv6.tar.gz)
- [Linux armv7](https://collection.b0.upaiyun.com/softwares/upx/upx_0.3.7_linux_armv7.tar.gz)

### 源码编译

> 需要安装 [Golang 编译环境](https://golang.org/dl/)

```
$ git clone https://github.com/upyun/upx.git
$ cd upx && make
```
or

```
$ GO111MODULE=on go get -u github.com/upyun/upx@v0.3.7
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

---

## 使用

> 所有命令都支持 `-h` 查看使用方法

|    命令  | 说明 |
| -------- | ---- |
| [login](#login)    | 登录又拍云存储 |
| [logout](#logout)   | 退出帐号 |
| [sessions](#sessions) | 查看所有的会话 |
| [switch](#switch)   | 切换会话 |
| [info](#info)     | 显示服务名、用户名等信息 |
| [ls](#ls)       | 显示当前目录下文件和目录信息 |
| [cd](#cd)       | 改变工作目录（进入一个目录）|
| [pwd](#pwd)      | 显示当前所在目录 |
| [mkdir](#mkdir)    | 创建目录 |
| [tree](#tree)     | 显示目录结构 |
| [get](#get)      | 下载一个文件或目录 |
| [put](#put)      | 上传一个文件或目录 |
| [upload](#upload)   | 上传多个文件或目录或 http(s) 文件, 支持 Glob 模式过滤上传文件|
| [rm](#rm)       | 删除目录或文件 |
| [sync](#sync)     | 目录增量同步，类似 rsync |
| [auth](#auth)     | 生成包含空间名操作员密码信息的 auth 字符串 |
| [post](#post)     | 提交异步处理任务 |
| [purge](#purge)    | 提交 CDN 缓存刷新任务 |


| global options | 说明 |
| -------------- | ---- |
| --quiet, -q    | 不显示信息 |
| --auth value   | auth 字符串 |
| --help, -h     | 显示帮助信息 |
| --version, -v  | 显示版本号 |


## login
> 使用又拍云操作员账号登录服务, 登录成功后将会保存会话，支持同时登录多个服务, 使用 `switch` 切换会话。

需要提供的验证字段
+ ServiceName: 服务(bucket)的名称
+ Operator: 操作员名
+ Password: 操作员密码


#### 语法
```bash
upx login
```

#### 示例
```bash
upx login

#ServiceName: testService
#Operator:  upx
#Password: password
```

## logout
> 退出当前登录的会话，如果存在多个登录的会话，可以使用 `switch` 切换到需要退出的会话，然后退出。


#### 语法
```bash
upx logout
```

#### 示例
```bash
upx logout

# Goodbye upx/testService ~
```


## sessions
> 列举出当前登录的所有会话

#### 语法
```bash
upx sessions
```

#### 示例
```bash
upx sessions

# > mybucket1
# > mybucket2
# > mybucket3
```

## switch
> 切换登录会话, 通过 `sessions` 命令可以查看所有的会话列表。

|  args  | 说明 |
| --------- | ---- |
| service-name | 服务名称(bucket) |

#### 语法
```bash
upx switch <service-name>
```

#### 示例
```bash
upx switch mybucket3
```

## info
> 查看当前服务的状态。

#### 语法
```bash
upx info
```

#### 示例
```bash
upx info
> ServiceName:   mybucket1
> Operator:      tester
> CurrentDir:    /
> Usage:         2.69GB
```

## ls
> 默认按文件修改时间先后顺序输出

|  args  | 说明 |
| --------- | ---- |
| remote-path | 远程路径 |

|  options  | 说明 |
| --------- | ---- |
| -d        | 仅显示目录 |
| -r        | 文件修改时间倒序输出 |
| --color   | 根据文件类型输出不同的颜色 |
| -c v      | 仅显示前 v 个文件或目录, 默认全部显示  |
| --mtime v | 通过文件被修改的时间删选，参考 Linux `find` |

#### 语法
```bash
upx ls [options...] [remote-path]
```

#### 示例

查看根目录下的文件
```bash
upx ls / 
```

只查看根目录下的目录
```bash
upx ls -d /
```

只查看根目录下的修改时间大于3天的文件
```bash
upx ls --mtime +3 /
```

只查看根目录下的修改时间小于1天的文件
```bash
upx ls --mtime -1 /
```

## cd
> 改变当前的工作路径，默认工作路径为根目录, 工作路径影响到操作时的默认远程路径。

|  args  | 说明 |
| --------- | ---- |
| remote-path | 远程路径 |

#### 语法
```bash
upx cd <remote-path>
```

#### 示例
将当前工作路径切换到 `/www`
```
upx cd /www
```

## pwd
> 显示当前所在的远程目录

#### 语法
```bash
upx pwd
```

#### 示例
```bash
upx pwd

> /www
```

## mkdir
> 创建远程目录

|  args  | 说明 |
| --------- | ---- |
| remote-dir | 远程目录 |

#### 语法
```bash
upx mkdir <remote-dir>
```

#### 示例
在当前工作目录下创建一个名为 mytestbucket 的目录
```bash
upx mkdir mytestbucket
```

在根目录下创建一个名为 mytestbucket 的目录
```bash
upx mkdir /mytestbucket
```

## tree
> 显示目录结构，树形模式显示

#### 语法
```bash
upx tree
```

#### 示例
查看 `/ccc` 目录下的目录结构
```bash
upx tree /ccc
> |-- aaacd
> !   |-- mail4788ca.png
> |-- ccc
> !   |-- Eroge de Subete wa Kaiketsu Dekiru! The Animation - 02 (2022) [1080p-HEVC-WEBRip][8D1929F5].mkv
> !   |-- baima_text_auditer.tar
> !   |-- linux-1.txt
```

## get
> 下载文件

|  args  | 说明 |
| --------- | ---- |
| remote-path | 远程路径，支持文件或文件夹 |
| saved-file | 需要保存到的本地目录，或指定完整的文件名 |

| options | 说明                          |
|---------|-----------------------------|
| -w | 多线程下载 (1-10) (default: 5) |
| -c   | 恢复中断的下载    |
| --start | 只下载路径字典序大于等于 `start` 的文件或目录 |
| --end   | 只下载路径字典序小于 `end` 的文件或目录     |


#### 语法
```bash
upx get [options] <remote-path> [saved-file]
```

#### 示例
下载文件
```bash
upx get /baima_text_auditer.tar
```

下载文件时指定保存路径
```bash
upx get /baima_text_auditer.tar ./baima_text_auditer2.tar
```

多线程下载文件
```bash
upx get -w 10 /baima_text_auditer.tar
```

恢复中断的下载
```bash
upx get -c /baima_text_auditer.tar
```

## put
> 上传文件或文件夹

|  args  | 说明 |
| --------- | ---- |
| local-file | 本地的文件或文件夹 |
| remote-file | 需要保存到的远程文件路径或文件夹 |

| options | 说明                          |
|---------|-----------------------------|
| -w | 多线程下载 (1-10) (default: 5) |

#### 语法
```bash
upx put <local-file> [remote-file]
```

#### 示例
上传本地文件，到远程绝对路径
```bash
upx put aaa.mp4 /video/aaa.mp4
```

上传本地目录，到远程绝对路径
```bash
upx put ./video /myfiles
```

## upload
> 上传文件或目录或 url 链接，支持多文件，文件名匹配

|  args  | 说明 |
| --------- | ---- |
| local-file | 本地的文件或文件夹, 或匹配文件规则 |
| remote-path | 需要保存到的远程文件路径 |

| options | 说明                          |
|---------|-----------------------------|
| -w | 多线程下载 (1-10) (default: 5) |
| --remote | 远程路径 |

#### 语法
```bash
upx upload [--remote remote-path] <local-file>
```

#### 示例

上传当前路径下的所有 `jpg` 图片到 `/images` 目录
```
upx upload --remote /images ./*.jpg
```

上传 `http` 文件到 `/files` 目录
```
upx upload --remote /files https://xxxx.com/myfile.tar.gz
```

## rm

> 默认不会删除目录，支持通配符 `*`

|  args  | 说明 |
| --------- | ---- |
| remote-file | 远程文件 |

|  options  | 说明 |
| --------- | ---- |
| -d        | 仅删除目录 |
| -a        | 删除目录跟文件 |
| --async   | 异步删除，目录可能需要二次删除 |
| --mtime v | 参考 Linux `find` |

#### 语法
```bash
upx rm [options] <remote-file>
```

#### 示例
删除目录 `/www`
```bash
upx rm -d /www
```

删除文件 `/aaa.png`
```bash
upx rm /aaa.png
```

## sync

> sync 本地路径 存储路径

|  args  | 说明 |
| --------- | ---- |
| local-path | 本地的路径 |
| remote-path | 远程文件路径 |

| options  | 说明 |
| -------- | ---- |
| -w       | 指定并发数，默认为 5 |
| --delete | 删除上一次同步后本地删除的文件 |

#### 语法
```bash
upx sync <local-path> <remote-path>
```

#### 示例
同步本地路径和远程路径
```bash
upx ./workspace /workspace
```

## auth

> 生成包含空间名操作员密码信息, auth 空间名 操作员 密码

#### 示例
当命令中包含 `--auth` 参数时，会忽略已登陆的信息。

```bash
upx auth mybucket user password
```


## post

|     options    | 说明 |
| -------------- | ---- |
| --app value    | app 名称 |
| --notify value | 回调地址 |
| --task value   | 任务文件名 |

## purge

> purge url --list urls

|     options    | 说明 |
| -------------- | ---- |
| --list value   | 批量刷新文件名 |


## TODO

- [x] put 支持断点续传
- [ ] upx 支持指定 API 地址
