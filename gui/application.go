package gui

import "github.com/digitalautonomy/keymirror/api"

type application struct {
	ui   *ui
	keys api.KeyAccess
}

