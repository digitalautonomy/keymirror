package gui

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type guiSuite struct {
	suite.Suite
}

func TestGUISuite(t *testing.T) {
	suite.Run(t, new(guiSuite))
}

func (s *guiSuite) Test_Start_StartsGTKApplication() {
	defer gostub.Stub(&glibi.APPLICATION_FLAGS_NONE, glibi.ApplicationFlags(77)).Reset()

	appMock := &gtk.MockApplication{}
	appMock.On("Connect", "activate", mock.AnythingOfType("func()")).Return(0).Once()
	appMock.On("Run", []string{}).Return(0).Once()

	gtkMock := &gtk.Mock{}
	gtkMock.On("ApplicationNew", "digital.autonomia.keymirror", glibi.APPLICATION_FLAGS_NONE).Return(appMock, nil).Once()

	Start(gtkMock)

	appMock.AssertExpectations(s.T())
	gtkMock.AssertExpectations(s.T())
}

func (s *guiSuite) Test_Start_ConnectsAnEventHandlerForActivateSignalThatShowsTheMainApplicationWindow() {
	appMock := &gtk.MockApplication{}
	var activateEventHandler func()

	appMock.On("Connect", "activate", mock.AnythingOfType("func()")).
		Return(0).
		Run(func(args mock.Arguments) {
			activateEventHandler = args.Get(1).(func())
		})
	appMock.On("Run", mock.Anything).Return(0)

	gtkMock := &gtk.Mock{}
	gtkMock.On("ApplicationNew", mock.Anything, mock.Anything).Return(appMock, nil)

	Start(gtkMock)

	winMock := &gtk.MockApplicationWindow{}
	winMock.On("ShowAll").Return().Once()
	gtkMock.On("ApplicationWindowNew", appMock).Return(winMock, nil).Once()

	activateEventHandler()

	gtkMock.AssertExpectations(s.T())
	winMock.AssertExpectations(s.T())
}
