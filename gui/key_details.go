package gui

import "github.com/coyim/gotk3adapter/gtki"

func (u *ui) populateKeyDetails(box gtki.Box) {
	b, _ := buildObjectFrom[gtki.Box](u, "KeyDetails")
	box.Add(b)
}
