VER= $(shell cat VERSION)

all:
	- go get -u
	go build -o upx -ldflags "-X main.version=$(VER)" .

release:
	GOOS=linux  GOARCH=amd64 go build -o upx-linux-amd64-$(VER) -ldflags "-X main.version=$(VER)" .
	GOOS=linux  GOARCH=386  go build -o upx-linux-i386-$(VER) -ldflags "-X main.version=$(VER)" .
	GOOS=darwin GOARCH=amd64 go build -o upx-darwin-amd64-$(VER) -ldflags "-X main.version=$(VER)" .

upload:
	./upx pwd
	./upx put upx-linux-amd64-$(VER) /softwares/upx/
	./upx put upx-linux-i386-$(VER)  /softwares/upx/
	./upx put upx-darwin-amd64-$(VER) /softwares/upx/

install:
	install -c upx /usr/local/bin
