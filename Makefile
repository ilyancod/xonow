EXECUTABLE_NAME := Xonow

BUILD_DIR := build
APP_DIR := cmd

EXECUTABLE_PATH := $(BUILD_DIR)/$(EXECUTABLE_NAME)

SOURCE_FILES := $(APP_DIR)/main.go

default: run

clean:
	rm -rf $(BUILD_DIR)

run: build_cli
	$(EXECUTABLE_PATH)

build_cli:
	@mkdir -p $(BUILD_DIR)
	go build -o $(EXECUTABLE_PATH) $(SOURCE_FILES)