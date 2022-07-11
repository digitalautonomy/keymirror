//go:build binary

package main

import (
	"github.com/coyim/gotk3adapter/gdka"
	"github.com/coyim/gotk3adapter/gioa"
	"github.com/coyim/gotk3adapter/gtka"
)

func init() {
	realGTK = gtka.Real
	realGDK = gdka.Real
	realGIO = gioa.Real
}
