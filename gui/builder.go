package gui

import (
	"embed"
	"fmt"
	"github.com/coyim/gotk3adapter/glibi"
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

func resourceDefinitionPath(name string) string {
	return definitionPath("resources", name, "gresource")
}

type builder struct {
	gtki.Builder
}

func (b *builder) get(name string) glibi.Object {
	o, _ := b.GetObject(name)
	return o
}

func (u *ui) builderFrom(name string) *builder {
	pathOfFile := interfaceDefinitionPath(name)
	b, _ := u.gtk.BuilderNew()
	content, _ := fs.ReadFile(getDefinitions(), pathOfFile)
	b.AddFromString(string(content))
	return &builder{b}
}

func buildObjectFrom[T any](u *ui, name string) (T, *builder) {
	b := u.builderFrom(name)
	return b.get(name).(T), b
}
