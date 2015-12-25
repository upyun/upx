> upx is a tool for managing files in UPYUN. ONLY support \*nix and darwin

## Feature Summary

- upload file and directory
- download file and directory
- remove file and directory
- make directory
- list directory


## Installation

```
$ git clone https://github.com/polym/upx.git
$ cd upx && go build .
```

## Usage

```
Usage:

	upx command [arguments]

The commands are:

	cd       Change working directory
	get      Get directory or file from UPYUN
	login    Log in UPYUN with username, password, bucket
	logout   Log out UPYUN
	ls       List directory or file
	mkdir    Make directory
	put      Put directory or file to UPYUN
	pwd      Print working directory
	rm       Remove one or more directories and files
	version  Print version

```


## TODO

- support for removing all files which are too old
