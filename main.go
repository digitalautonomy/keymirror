package main

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/gui"
	"github.com/digitalautonomy/keymirror/ssh"
)

var realGTK gtki.Gtk = nil
var realGDK gdki.Gdk = nil
var startGUI = gui.Start

func main() {
	startGUI(realGTK, realGDK, ssh.Access)
}
