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
	goreleaser --rm-dist

upload: release
	./upx pwd
	./upx put dist/upx_darwin_amd64/upx /softwares/upx/upx-darwin-amd64-$(VER); \

	for ARCH in amd64 386 arm64 arm_6 arm_7; do \
		./upx put dist/upx_linux_$$ARCH/upx /softwares/upx/upx-linux-$$ARCH-$(VER); \
	done

	for ARCH in amd64 386; do \
		./upx put dist/upx_windows_$$ARCH/upx.exe /softwares/upx/upx-windows-$$ARCH-$(VER).exe; \
	done

.PHONY: app test release upload
