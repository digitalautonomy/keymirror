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
DEFINITION_DIR := gui/definitions
ICONS_RESOURCE_FILE := $(DEFINITION_DIR)/resources/icons.gresource
INTERFACE_DEFINITION_FILES := $(DEFINITION_DIR)/interface/*.xml
STYLES_DEFINITION_FILES := $(DEFINITION_DIR)/styles/*.css
RESOURCES_DEFINITION_FILES := $(ICONS_RESOURCE_FILE)
DEFINITION_FILES := $(INTERFACE_DEFINITION_FILES) $(STYLES_DEFINITION_FILES) $(RESOURCES_DEFINITION_FILES)
SOURCE_FILES := $(GO_FILES)

GO := go
GOBUILD := $(GO) build
GOTEST := $(GO) test
GOINSTALL := $(GO) install

COVERPROFILE := coverprofile

.PHONY := default clean test ci-test ci-deps ci-upload-coverage resources

default: $(BINARY)

$(BUILD_DIR):
	mkdir -p $@

$(BINARY): $(BUILD_DIR) $(SOURCE_FILES) $(DEFINITION_FILES) gui/definitions/resources/icons.gresource
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

$(ICONS_RESOURCE_FILE): resources/icons.xml resources/icons/*.png
	glib-compile-resources --target=$@ resources/icons.xml

resources: $(RESOURCES_DEFINITION_FILES)