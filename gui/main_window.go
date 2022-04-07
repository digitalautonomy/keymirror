package gui

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/api"
)

const keymirrorApplicationID = "digital.autonomia.keymirror"

func (a *application) createMainWindow(app gtki.Application) gtki.Window {
	w, b := buildObjectFrom[gtki.ApplicationWindow](a.ui, "MainWindow")
	box := b.get("keyListBox").(gtki.Box)
	a.ui.populateListWithKeyEntries(a.keys, box)
	w.SetApplication(app)
	return w
}

func (a *application) activate(app gtki.Application) {
	a.ui.applyApplicationStyle()
	mainWindow := a.createMainWindow(app)
	mainWindow.ShowAll()
}

func (a *application) start() {
	app, _ := a.ui.gtk.ApplicationNew(keymirrorApplicationID, glibi.APPLICATION_FLAGS_NONE)
	app.Connect("activate", func() { a.activate(app) })
	app.Run([]string{})
}

func Start(gtk gtki.Gtk, gdk gdki.Gdk, ka api.KeyAccess) {
	app := &application{
		ui: &ui{
			gtk: gtk,
			gdk: gdk,
		},

		keys: ka,
	}

	app.start()
}
