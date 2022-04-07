package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/api"
)

func (u *ui) createKeyEntryBoxFrom(entry api.KeyEntry) gtki.Box {
	b, _ := buildObjectFrom[gtki.Box](u, "KeyListEntry")
	b.GetChildren()[0].(gtki.Label).SetLabel(entry.Locations()[0])
	return b
}

func (u *ui) populateListWithKeyEntries(access api.KeyAccess, box gtki.Box) {
	for _, e := range access.AllKeys() {
		box.Add(u.createKeyEntryBoxFrom(e))
	}
}
