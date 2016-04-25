> upx is a tool for managing files in UPYUN. ONLY support \*nix and darwin

[![Build Status](https://travis-ci.org/polym/upx.svg?branch=master)](https://travis-ci.org/polym/upx)

## Feature Summary

- upload file and folder
- download file and folder
- remove file and folder, support wildcard, like `upx rm *jpg`
- make directory
- list directory
- support progress bar
- support multi-users


## Installation

### Download binary

```
$ wget -O /usr/local/bin/upx http://collection.b0.upaiyun.com/softwares/upx/upx-darwin-amd64-v0.1.1
$ wget -O /usr/local/bin/upx http://collection.b0.upaiyun.com/softwares/upx/upx-linux-amd64-v0.1.1
$ wget -O /usr/local/bin/upx http://collection.b0.upaiyun.com/softwares/upx/upx-linux-i386-v0.1.1
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
    cd          Change working directory
    get         Get directory or file from UPYUN
    info        Current information
    login       Log in UPYUN with service_name, username, password
    logout      Log out UPYUN
    ls          List directory or file
    mkdir       Make directory
    put         Put directory or file to UPYUN
    pwd         Print working directory
    rm          Remove one or more directories and files
    services    List all services
    switch      Switch service

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version
```


## TODO

- sync local files to UPYUN
- support for removing all files which are too old
- more options for commands
