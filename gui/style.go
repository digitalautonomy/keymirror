package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"io/fs"
)

func (u *ui) createStyleProviderFrom(filename string) gtki.CssProvider {
	cssProvider, _ := u.gtk.CssProviderNew()
	pathOfFile := styleDefinitionPath(filename)
	content, _ := fs.ReadFile(getDefinitions(), pathOfFile)
	cssProvider.LoadFromData(string(content))
	return cssProvider
}

func (u *ui) applyApplicationStyle() {
	s, _ := u.gdk.ScreenGetDefault()

	globalStyleProvider := u.createStyleProviderFrom("global")
	colorsStyleProvider := u.createStyleProviderFrom("colors_light")

	u.gtk.AddProviderForScreen(s, globalStyleProvider, uint(gtki.STYLE_PROVIDER_PRIORITY_APPLICATION))
	u.gtk.AddProviderForScreen(s, colorsStyleProvider, uint(gtki.STYLE_PROVIDER_PRIORITY_APPLICATION))
}
