> upx is a tool for managing files in UPYUN. Mac, Linux, Windows supported

[![Build Status](https://travis-ci.org/polym/upx.svg?branch=master)](https://travis-ci.org/polym/upx)

## Feature Summary

- [x] basic filesystem commands, like `mkdir`, `cd`, `ls`, `rm`, `pwd`
- [x] upload items to UPYUN
- [x] download items from UPYUN
- [x] rsync directory to UPYUN
- [x] remove by wildcard, like `upx rm *jpg`
- [x] a multi-user tool


## Installation

### Download binary

```
// mac
$ wget -O /usr/local/bin/upx http://collection.b0.upaiyun.com/softwares/upx/upx-darwin-amd64-v0.1.2

// linux
$ wget -O /usr/local/bin/upx http://collection.b0.upaiyun.com/softwares/upx/upx-linux-amd64-v0.1.2
$ wget -O /usr/local/bin/upx http://collection.b0.upaiyun.com/softwares/upx/upx-linux-i386-v0.1.2

// windows
$ wget -O /usr/local/bin/upx http://collection.b0.upaiyun.com/softwares/upx/upx-windows-amd64-v0.1.2.exe
$ wget -O /usr/local/bin/upx http://collection.b0.upaiyun.com/softwares/upx/upx-windows-i386-v0.1.2.exe

// mac or linux need
$ chmod +x /usr/local/bin/upx
```

### Source Compile

```
$ git clone https://github.com/polym/upx.git
$ cd upx && make
```

or

```
$ go get github.com/polym/upx
```

## Usage

```
NAME:
   upx - a tool for managing files in UPYUN

USAGE:
   upx [global options] command [command options] [arguments...]

COMMANDS:
     cd                 Change working directory
     get                Get directory or file from UPYUN
     info, i            Current information
     login              Log in UPYUN with service_name, username, password
     logout             Log out UPYUN
     ls                 List directory or file
     mkdir, mk          Make directory
     put                Put directory or file to UPYUN
     pwd                Print working directory
     rm                 Remove one or more directories and files
     services, sv       List all services, alias sv
     switch, sw         Switch service, alias sw
     sync               sync folder to UPYUN

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version
```
