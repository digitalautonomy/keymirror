
BUILD_DIR := build
BINARY := $(BUILD_DIR)/keymirror

GO_FILES := *.go
SOURCE_FILES := $(GO_FILES)

.PHONY := default

default: $(BINARY)

$(BUILD_DIR):
	mkdir -p $@

$(BINARY): $(BUILD_DIR) $(SOURCE_FILES)
	go build -o $@
