ifndef VER
	VER= 'latest'
endif

APP= upx
ROOT= $(shell echo $(GOPATH) | awk -F':' '{print $$1}')
PROJ_DIR= $(ROOT)/src/upyun.com
PWD= $(shell pwd)

app:
	go build -o $(APP) .

test:
	go test -v .

release:
	rm -rf release
	for OS in linux windows darwin; do \
		for ARCH in amd64 386; do \
			GOOS=$$OS GOARCH=$$ARCH go build -o release/upx-$$OS-$$ARCH-$(VER) .; \
		done \
	done
	tar -zcf release/upx-$(VER).tar.gz release/*

upload: release
	./upx pwd
	for OS in linux darwin; do \
		for ARCH in amd64 386; do \
			./upx put release/upx-$$OS-$$ARCH-$(VER) /softwares/upx/; \
		done \
	done
	for ARCH in amd64 386; do \
		./upx put release/upx-windows-$$ARCH-$(VER) /softwares/upx/upx-windows-$$ARCH-$(VER).exe; \
	done

.PHONY: app test release upload
