package main

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3mocks/gdk"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/digitalautonomy/keymirror/api"
	"github.com/prashantv/gostub"
	"github.com/sirupsen/logrus"
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
	var calledWithLog logrus.Ext1FieldLogger
	var calledWithKeyAccesss api.KeyAccess
	defer gostub.Stub(&startGUI, func(g gtki.Gtk, g2 gdki.Gdk, log logrus.Ext1FieldLogger, ka api.KeyAccess) {
		calledWithGTK = g
		calledWithGDK = g2
		calledWithLog = log
		calledWithKeyAccesss = ka
	}).Reset()

	main()

	s.Equal(ourGTK, calledWithGTK)
	s.Equal(ourGDK, calledWithGDK)
	s.NotNil(calledWithLog)
	s.Equal(logrus.TraceLevel, calledWithLog.(*logrus.Logger).Level)
	s.NotNil(calledWithKeyAccesss)
}
