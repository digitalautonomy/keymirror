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

func buildObjectFrom[T any](gtk gtki.Gtk, name string) T {
	pathOfFile := fmt.Sprintf("definitions/interface/%s.xml", name)
	builder, _ := gtk.BuilderNew()
	content, _ := fs.ReadFile(getDefinitions(), pathOfFile)
	builder.AddFromString(string(content))
	w, _ := builder.GetObject(name)

	return w.(T)
}
