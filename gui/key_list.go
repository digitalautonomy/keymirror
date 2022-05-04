package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/api"
	"github.com/digitalautonomy/keymirror/i18n"
)

func (u *ui) createKeyEntryBoxFrom(entry api.KeyEntry) gtki.Box {
	b, _ := buildObjectFrom[gtki.Box](u, "KeyListEntry")
	b.GetChildren()[0].(gtki.Label).SetLabel(entry.Locations()[0])
	return b
}

func (u *ui) populateListWithKeyEntries(access api.KeyAccess, box gtki.Box, onNoKeys func(box gtki.Box)) {
	for _, e := range access.AllKeys() {
		onNoKeys = func(box gtki.Box) {}
		box.Add(u.createKeyEntryBoxFrom(e))
	}
	onNoKeys(box)
}

func (u *ui) showNoAvailableKeysMessage(box gtki.Box) {
	l, _ := u.gtk.LabelNew(i18n.Local("\u26A0 No keys available \u26A0"))
	sc, _ := l.GetStyleContext()
	sc.AddClass("infoMessage")
	box.Add(l)
}
