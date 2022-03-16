
BUILD_DIR := build
BINARY := $(BUILD_DIR)/keymirror

GO_FILES := *.go
SOURCE_FILES := $(GO_FILES)

GO := go
GOBUILD := $(GO) build
GOTEST := $(GO) test
GOINSTALL := $(GO) install

COVERPROFILE := coverprofile

.PHONY := default clean test ci-test ci-deps ci-upload-coverage

default: $(BINARY)

$(BUILD_DIR):
	mkdir -p $@

$(BINARY): $(BUILD_DIR) $(SOURCE_FILES)
	$(GOBUILD) -o $@

clean:
	$(RM) -r $(BUILD_DIR)

test:
	$(GOTEST) ./...

coverage:
	$(GOTEST) -cover -coverprofile coverlog ./... || true
	$(GO) tool cover -html coverlog
	$(RM) coverlog

$(COVERPROFILE):
	$(GOTEST) -cover -coverprofile $@ ./...

ci-test: $(COVERPROFILE)

ci-deps:
	$(GOINSTALL) github.com/mattn/goveralls@latest

ci-upload-coverage: $(COVERPROFILE) ci-deps
	goveralls -coverprofile=$(COVERPROFILE)
