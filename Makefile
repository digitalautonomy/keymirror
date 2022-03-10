
BUILD_DIR := build
BINARY := $(BUILD_DIR)/keymirror

GO_FILES := *.go
SOURCE_FILES := $(GO_FILES)

GO := go
GOBUILD := $(GO) build
GOTEST := $(GO) test


.PHONY := default clean test

default: $(BINARY)

$(BUILD_DIR):
	mkdir -p $@

$(BINARY): $(BUILD_DIR) $(SOURCE_FILES)
	$(GOBUILD) -o $@

clean: 
	$(RM) -r $(BUILD_DIR)

test:
	$(GOTEST)

coverage:
	$(GOTEST) -cover -coverprofile coverlog || true
	$(GO) tool cover -html coverlog
	$(RM) coverlog