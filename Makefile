EXECUTABLE_NAME := Xonow

BUN_DIR := bin
APP_DIR := cmd

EXECUTABLE_PATH := $(BUN_DIR)/$(EXECUTABLE_NAME)

SOURCE_FILES := $(APP_DIR)/main.go

SYSTEMD_UNIT=xonow.service

default: deploy

clean:
	rm -rf $(BUN_DIR)

run: build_cli
	$(EXECUTABLE_PATH)

build_cli:
	@mkdir -p $(BUN_DIR)
	go build -o $(EXECUTABLE_PATH) $(SOURCE_FILES)

restart:
	@systemctl --user restart $(SYSTEMD_UNIT)

deploy: build_cli restart