package gui

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

func Start(gtk gtki.Gtk) {
	app, _ := gtk.ApplicationNew("digital.autonomia.keymirror", glibi.APPLICATION_FLAGS_NONE)
	app.Connect("activate", func() {
		w, _ := gtk.ApplicationWindowNew(app)
		w.ShowAll()
	})

	app.Run([]string{})
}
