GTK_VERSION := $(shell pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
GTK_VERSION_TAG := "gtk_$(GTK_VERSION)"

GLIB_VERSION := $(shell pkg-config --modversion glib-2.0 | tr . _ | cut -d '_' -f 1-2)
GLIB_VERSION_TAG := "glib_$(GLIB_VERSION)"

GDK_VERSION := $(shell pkg-config --modversion gdk-3.0 | tr . _ | cut -d '_' -f 1-2)
GDK_VERSION_TAG := "gdk_$(GDK_VERSION)"

BINARY_TAGS := -tags $(GTK_VERSION_TAG),$(GLIB_VERSION_TAG),$(GDK_VERSION_TAG),binary

BUILD_DIR := build
BINARY := $(BUILD_DIR)/keymirror

GO_FILES := *.go ssh/*.go gui/*.go
INTERFACE_DEFINITION_FILES := gui/definitions/interface/*.xml
STYLES_DEFINITION_FILES := gui/definitions/styles/*.css
DEFINITION_FILES := $(INTERFACE_DEFINITION_FILES) $(STYLES_DEFINITION_FILES)
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

$(BINARY): $(BUILD_DIR) $(SOURCE_FILES) $(DEFINITION_FILES)
	$(GOBUILD) $(BINARY_TAGS) -o $@

clean:
	$(RM) -r $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

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

coverage-tails: 
	$(GOTEST) -cover -coverprofile coverlog ./... || true
	$(GO) tool cover -html coverlog -o ~/Tor\ Browser/coverage.html
	xdg-open ~/Tor\ Browser/coverage.html
	$(RM) coverlog
