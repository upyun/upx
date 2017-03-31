ifndef VER
	VER= 'latest'
endif

APP= upx
ROOT= $(shell echo $(GOPATH) | awk -F':' '{print $$1}')
PROJ_DIR= $(ROOT)/src/upyun.com
PWD= $(shell pwd)

app:
	- mkdir -p $(PROJ_DIR) && ln -s $(PWD) $(PROJ_DIR)/$(APP)
	cd $(PROJ_DIR)/$(APP) && go build -o $(APP) .
	unlink $(PROJ_DIR)/$(APP)

vendor:
	- mkdir -p $(PROJ_DIR) && ln -s $(PWD) $(PROJ_DIR)/$(APP)
	cd $(PROJ_DIR)/$(APP) && govendor init && govendor add +external
	unlink $(PROJ_DIR)/$(APP)

test:
	- mkdir -p $(PROJ_DIR) && ln -s $(PWD) $(PROJ_DIR)/$(APP)
	cd $(PROJ_DIR)/$(APP) && go test -v .
	unlink $(PROJ_DIR)/$(APP)

release:
	- mkdir -p $(PROJ_DIR) && ln -s $(PWD) $(PROJ_DIR)/$(APP)
	cd $(PROJ_DIR)/$(APP) && for OS in linux windows darwin; do \
		for ARCH in amd64 386; do \
			GOOS=$$OS GOARCH=$$ARCH go build -o upx-$$OS-$$ARCH-$(VER) .; \
			GOOS=$$OS GOARCH=$$ARCH go test -c -o upx-$$OS-$$ARCH-$(VER).test .; \
		done \
	done
	unlink $(PROJ_DIR)/$(APP)

upload: release
	./upx pwd
	for OS in linux darwin; do \
		for ARCH in amd64 386; do \
			./upx put upx-$$OS-$$ARCH-$(VER) /softwares/upx/; \
		done \
	done
	for ARCH in amd64 386; do \
		./upx put upx-$$OS-$$ARCH-$(VER) /softwares/upx/upx-$$OS-$$ARCH-$(VER).exe; \
	done

.PHONY: app vendor test release upload
