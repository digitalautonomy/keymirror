package gui

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

func Start(gtk gtki.Gtk, gdk gdki.Gdk) {
	u := &ui{
		gtk: gtk,
		gdk: gdk,
	}

	app, _ := u.gtk.ApplicationNew("digital.autonomia.keymirror", glibi.APPLICATION_FLAGS_NONE)
	app.Connect("activate", func() {
		u.applyApplicationStyle()
		w := buildObjectFrom[gtki.ApplicationWindow](u, "MainWindow")
		w.SetApplication(app)
		w.ShowAll()
	})

	app.Run([]string{})
}
