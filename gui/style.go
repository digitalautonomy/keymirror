package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"io/fs"
)

func createStyleProviderFrom(gtk gtki.Gtk, filename string) gtki.CssProvider {
	cssProvider, _ := gtk.CssProviderNew()
	pathOfFile := styleDefinitionPath(filename)
	content, _ := fs.ReadFile(getDefinitions(), pathOfFile)
	cssProvider.LoadFromData(string(content))
	return cssProvider
}
