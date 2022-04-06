package gui

import (
	"github.com/coyim/gotk3adapter/gdki"
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

func applyApplicationStyle(gtk gtki.Gtk, gdk gdki.Gdk) {
	s, _ := gdk.ScreenGetDefault()

	globalStyleProvider := createStyleProviderFrom(gtk, "global")
	colorsStyleProvider := createStyleProviderFrom(gtk, "colors_light")

	gtk.AddProviderForScreen(s, globalStyleProvider, uint(gtki.STYLE_PROVIDER_PRIORITY_APPLICATION))
	gtk.AddProviderForScreen(s, colorsStyleProvider, uint(gtki.STYLE_PROVIDER_PRIORITY_APPLICATION))
}
