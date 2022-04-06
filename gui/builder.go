package gui

import (
	"embed"
	"fmt"
	"github.com/coyim/gotk3adapter/gtki"
	"io/fs"
)

//go:embed definitions
var definitionsInternal embed.FS

var getDefinitions = func() fs.FS {
	return definitionsInternal
}

func definitionPath(definitionType, name, extension string) string {
	return fmt.Sprintf("definitions/%s/%s.%s", definitionType, name, extension)
}

func interfaceDefinitionPath(name string) string {
	return definitionPath("interface", name, "xml")
}

func styleDefinitionPath(name string) string {
	return definitionPath("styles", name, "css")
}

func buildObjectFrom[T any](gtk gtki.Gtk, name string) T {
	pathOfFile := interfaceDefinitionPath(name)
	builder, _ := gtk.BuilderNew()
	content, _ := fs.ReadFile(getDefinitions(), pathOfFile)
	builder.AddFromString(string(content))
	w, _ := builder.GetObject(name)

	return w.(T)
}
