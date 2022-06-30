package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/api"
)

type clearable[T any] interface {
	GetChildren() []T
	Remove(T)
}

func clearAllChildrenOf[T any](b clearable[T]) {
	for _, c := range b.GetChildren() {
		b.Remove(c)
	}
}

func (u *ui) populateKeyDetails(key api.KeyEntry, into gtki.Box) {
	clearAllChildrenOf[gtki.Widget](into)

	p := key.PublicKeyLocations()

	if p != nil {
		b, builder := buildObjectFrom[gtki.Box](u, "KeyDetails")
		label := builder.get("publicKeyPath").(gtki.Label)
		label.SetLabel(p[0])
		label.SetTooltipText(p[0])
		into.Add(b)
	}
}
