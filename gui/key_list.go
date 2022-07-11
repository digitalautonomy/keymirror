package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/api"
	"github.com/digitalautonomy/keymirror/i18n"
)

func (u *ui) createKeyEntryBoxFrom(entry api.KeyEntry, detailsBox gtki.Box, detailsRev gtki.Revealer) gtki.Widget {
	b, builder := buildObjectFrom[gtki.Button](u, "KeyListEntry")
	builder.get("keyListEntryLabel").(gtki.Label).SetLabel(entry.Locations()[0])
	b.Connect("clicked", func() {
		// We have several possibilities:
		// - no box is currently visible because we just started the program
		// - no box is visible because it was previously collapsed
		// - the box is visible with information about another key
		// - the box is visible with information about the same key as currently clicked
		u.populateKeyDetails(entry, detailsBox)

		if u.currentlyVisibleKeyEntryButton != nil {
			removeClass(*u.currentlyVisibleKeyEntryButton, "current")
		}

		if u.currentlyVisibleKeyEntry == nil || *u.currentlyVisibleKeyEntry != entry {
			detailsRev.Show()
			detailsRev.SetRevealChild(true)
			addClass(b, "current")
			u.currentlyVisibleKeyEntry = &entry
			u.currentlyVisibleKeyEntryButton = &b
		} else {
			detailsRev.SetRevealChild(false)
			detailsRev.Hide()
			u.currentlyVisibleKeyEntry = nil
			u.currentlyVisibleKeyEntryButton = nil
		}
		u.onWindowSizeChange()
	})
	return b
}

func (u *ui) populateListWithKeyEntries(access api.KeyAccess, box gtki.Box, detailsBox gtki.Box, detailsRev gtki.Revealer, onNoKeys func(box gtki.Box)) {
	for _, e := range access.AllKeys() {
		onNoKeys = func(box gtki.Box) {}
		box.Add(u.createKeyEntryBoxFrom(e, detailsBox, detailsRev))
	}
	onNoKeys(box)
}

func (u *ui) showNoAvailableKeysMessage(box gtki.Box) {
	l, _ := u.gtk.LabelNew(i18n.Local("\u26A0 No keys available \u26A0"))
	sc, _ := l.GetStyleContext()
	sc.AddClass("infoMessage")
	box.Add(l)
}
