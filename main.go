package main

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/gui"
)

var realGTK gtki.Gtk = nil
var startGUI = gui.Start

func main() {
	startGUI(realGTK)
}
