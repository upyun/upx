> upx is a tool for managing files in UPYUN. ONLY support \*nix and darwin

## Feature Summary

- upload file and folder
- download file and folder
- remove file and folder, support wildcard, like `upx rm *jpg`
- make directory
- list directory
- support progress bar
- support multi-users **NEW**


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
Usage:

        upx command [arguments]

The commands are:

        cd       Change working directory
        get      Get directory or file from UPYUN
        help     Help information
        info     Current information
        login    Log in UPYUN with username, password, bucket
        logout   Log out UPYUN
        ls       List directory or file
        mkdir    Make directory
        put      Put directory or file to UPYUN
        pwd      Print working directory
        rm       Remove one or more directories and files
        sevices  List all services
        switch   Switch service
        version  Print version

```


## TODO

- list large directory (only list first 200 now)
- sync local files to UPYUN
- support for removing all files which are too old
- more options for command
- windows support
