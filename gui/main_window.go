package gui

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

func Start(gtk gtki.Gtk, gdk gdki.Gdk) {
	app, _ := gtk.ApplicationNew("digital.autonomia.keymirror", glibi.APPLICATION_FLAGS_NONE)
	app.Connect("activate", func() {
		applyApplicationStyle(gtk, gdk)
		w := buildObjectFrom[gtki.ApplicationWindow](gtk, "MainWindow")
		w.SetApplication(app)
		w.ShowAll()
	})

	app.Run([]string{})
}
