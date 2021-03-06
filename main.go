package main

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gioi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/gui"
	"github.com/digitalautonomy/keymirror/ssh"
	"github.com/sirupsen/logrus"
)

var realGTK gtki.Gtk = nil
var realGDK gdki.Gdk = nil
var realGIO gioi.Gio = nil
var startGUI = gui.Start

func main() {
	l := logrus.New()
	l.Level = logrus.TraceLevel
	startGUI(realGTK, realGDK, realGIO, l, ssh.Access(l))
}
