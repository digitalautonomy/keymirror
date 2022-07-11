package gui

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gioi"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/api"
	"github.com/sirupsen/logrus"
)

const keymirrorApplicationID = "digital.autonomia.keymirror"
const keymirrorApplicationResourceID = "/digital/autonomia/KeyMirror"

func (a *application) createMainWindow(app gtki.Application) gtki.Window {
	w, b := buildObjectFrom[gtki.ApplicationWindow](a.ui, "MainWindow")
	box := b.get("keyListBox").(gtki.Box)
	box2 := b.get("keyDetailsBox").(gtki.Box)
	keyDetailsRevealer := b.get("keyDetailsRevealer").(gtki.Revealer)
	a.populateMainWindow(box, box2, keyDetailsRevealer)
	w.SetApplication(app)
	return w
}

func (a *application) populateMainWindow(listBox, detailsBox gtki.Box, detailsRev gtki.Revealer) {
	a.ui.populateListWithKeyEntries(a.keys, listBox, detailsBox, detailsRev, a.ui.showNoAvailableKeysMessage)
}

func (a *application) activate(app gtki.Application) {
	a.ui.loadResourceDefinitions()
	a.ui.applyApplicationStyle()
	mainWindow := a.createMainWindow(app)
	mainWindow.ShowAll()
	a.ui.onWindowSizeChange = func() {
		mainWindow.Resize(1, mainWindow.GetAllocatedHeight())
	}
	a.ui.onWindowSizeChange()
}

func (a *application) start() {
	app, _ := a.ui.gtk.ApplicationNew(keymirrorApplicationID, glibi.APPLICATION_FLAGS_NONE)
	app.Connect("activate", func() { a.activate(app) })
	app.Run([]string{})
}

func Start(gtk gtki.Gtk, gdk gdki.Gdk, gio gioi.Gio, log logrus.Ext1FieldLogger, ka api.KeyAccess) {
	app := &application{
		ui: &ui{
			gtk: gtk,
			gdk: gdk,
			gio: gio,
			log: log.WithField("component", "gui"),
		},

		keys: ka,
	}

	app.start()
}
