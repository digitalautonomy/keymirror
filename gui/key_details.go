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
	b, builder := buildObjectFrom[gtki.Box](u, "KeyDetails")
	builder.get("publicKeyPath").(gtki.Label).SetLabel(key.PublicKeyLocations()[0])
	clearAllChildrenOf[gtki.Widget](into)
	into.Add(b)
}
