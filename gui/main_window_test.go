package gui

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3mocks/gdk"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"testing/fstest"
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

	gdkMock := &gdk.Mock{}
	Start(gtkMock, gdkMock)

	appMock.AssertExpectations(s.T())
	gtkMock.AssertExpectations(s.T())
}

func setupStubbedDefinitions() func() {
	definitionsMock := fstest.MapFS{
		"definitions/interface/MainWindow.xml": &fstest.MapFile{
			Data: []byte("this is an interface description for the object"),
		},
		"definitions/styles/global.css": &fstest.MapFile{
			Data: []byte("global content"),
		},
		"definitions/styles/colors_light.css": &fstest.MapFile{
			Data: []byte("light colors content"),
		},
	}

	return gostub.StubFunc(&getDefinitions, definitionsMock).Reset
}

func stubStyleProviders(gtkMock *gtk.Mock, gdkMock *gdk.Mock) {
	cssProvider := &gtk.MockCssProvider{}
	screenMock := &gdk.MockScreen{}
	gtkMock.On("CssProviderNew").Return(cssProvider, nil)
	cssProvider.On("LoadFromData", mock.Anything).Return(nil)

	gtkMock.On("AddProviderForScreen", mock.Anything, mock.Anything, mock.Anything).Return()
	gdkMock.On("ScreenGetDefault").Return(screenMock, nil).Once()
}

func mockObjectBuild(gtkMock *gtk.Mock, objectName string, ret interface{}) {
	builderMock := &gtk.MockBuilder{}

	fileContent := "this is an interface description for the object"

	gtkMock.On("BuilderNew").Return(builderMock, nil)
	builderMock.On("AddFromString", fileContent).Return(nil)
	builderMock.On("GetObject", objectName).Return(ret, nil)
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

	gdkMock := &gdk.Mock{}
	Start(gtkMock, gdkMock)

	winMock := &gtk.MockApplicationWindow{}
	winMock.On("SetApplication", appMock).Return().Once()
	winMock.On("ShowAll").Return().Once()

	mockObjectBuild(gtkMock, "MainWindow", winMock)
	stubStyleProviders(gtkMock, gdkMock)
	defer setupStubbedDefinitions()()

	activateEventHandler()

	gtkMock.AssertExpectations(s.T())
	winMock.AssertExpectations(s.T())
}
