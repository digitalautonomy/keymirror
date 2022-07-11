package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"io/fs"
)

//
//// #cgo pkg-config: gdk-3.0 gio-2.0 glib-2.0 gobject-2.0 gtk+-3.0
//// #include <gio/gio.h>
//// #include <gtk/gtk.h>
//import "C"
//
//func withDefinitionInTempFile(def string, f func(string)) {
//	content, _ := fs.ReadFile(getDefinitions(), def)
//	td, _ := ioutil.TempDir("", "")
//	defer func() {
//		_ = os.RemoveAll(td)
//	}()
//	file := path.Join(td, path.Base(def))
//	_ = ioutil.WriteFile(file, content, 0755)
//	f(file)
//}
//
//func testBla() {
//	withDefinitionInTempFile("definitions/icons.gresource", func(path string) {
//		gr, _ := gio.LoadGResource(path)
//		gio.RegisterGResource(gr)
//		C.gtk_icon_theme_add_resource_path(C.gtk_icon_theme_get_default(), C.CString("/digital/autonomia/KeyMirror"))
//	})
//}

func (u *ui) createStyleProviderFrom(filename string) gtki.CssProvider {
	cssProvider, _ := u.gtk.CssProviderNew()
	pathOfFile := styleDefinitionPath(filename)
	u.log.WithField("file", pathOfFile).Debug("loading CSS style")
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
