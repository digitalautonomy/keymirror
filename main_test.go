package main

import (
	"github.com/coyim/gotk3adapter/gtki"
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

	var calledWithG gtki.Gtk
	defer gostub.Stub(&startGUI, func(g gtki.Gtk) {
		calledWithG = g
	}).Reset()

	main()

	s.Equal(ourGTK, calledWithG)
}
