package gui

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gioi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/api"
	"github.com/sirupsen/logrus"
)

// ui contains the core things we need for the UI to function
// it will contain the implementations for glib, gtk, pango and gdk,
// logging and error handling. most other things does NOT belong here.

type ui struct {
	gtk gtki.Gtk
	gdk gdki.Gdk
	gio gioi.Gio

	// error handler
	log logrus.Ext1FieldLogger

	currentlyVisibleKeyEntry       *api.KeyEntry
	currentlyVisibleKeyEntryButton *gtki.Button
	onWindowSizeChange             func()
}
