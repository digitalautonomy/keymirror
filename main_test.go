package main

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3mocks/gdk"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/prashantv/gostub"
	"testing"

	"github.com/stretchr/testify/suite"
)

type mainSuite struct {
	suite.Suite
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(mainSuite))
}

func (s *mainSuite) Test_main_startsTheGuiWithTheRealGTK() {
	originalGTK := realGTK
	defer func() {
		realGTK = originalGTK
	}()
	ourGTK := &gtk.Mock{}
	realGTK = ourGTK

	originalGDK := realGDK
	defer func() {
		realGDK = originalGDK
	}()
	ourGDK := &gdk.Mock{}
	realGDK = ourGDK

	var calledWithGTK gtki.Gtk
	var calledWithGDK gdki.Gdk
	defer gostub.Stub(&startGUI, func(g gtki.Gtk, g2 gdki.Gdk) {
		calledWithGTK = g
		calledWithGDK = g2
	}).Reset()

	main()

	s.Equal(ourGTK, calledWithGTK)
	s.Equal(ourGDK, calledWithGDK)
}
