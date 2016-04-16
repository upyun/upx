VER= $(shell cat VERSION)

all:
	go get -d
	go build -o upx

release:
	GOOS=linux  GOARCH=amd64 go build -o upx-linux-amd64-$(VER) .
	GOOS=linux  GOARCH=386  go build -o upx-linux-i386-$(VER) .
	GOOS=darwin GOARCH=amd64 go build -o upx-darwin-amd64-$(VER) .

upload:
	./upx pwd
	./upx put upx-linux-amd64-$(VER) /softwares/upx/
	./upx put upx-linux-i386-$(VER)  /softwares/upx/
	./upx put upx-darwin-amd64-$(VER) /softwares/upx/

install:
	install -c upx /usr/local/bin

test:
	go test -v
