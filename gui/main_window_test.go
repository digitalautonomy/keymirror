package gui

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3mocks/gdk"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/prashantv/gostub"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"testing/fstest"
)

type expectationsAsserter interface {
	AssertExpectations(mock.TestingT) bool
}

type guiSuite struct {
	suite.Suite

	gtkMock         *gtk.Mock
	objectsToAssert []expectationsAsserter
}

func (s *guiSuite) addObjectToAssert(o expectationsAsserter) {
	s.objectsToAssert = append(s.objectsToAssert, o)
}

func (s *guiSuite) SetupTest() {
	s.objectsToAssert = []expectationsAsserter{}
	s.gtkMock = &gtk.Mock{}
	s.addObjectToAssert(s.gtkMock)
}

func (s *guiSuite) TearDownTest() {
	for _, t := range s.objectsToAssert {
		t.AssertExpectations(s.T())
	}
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
	log, _ := test.NewNullLogger()
	Start(gtkMock, gdkMock, log, nil)

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

func mockObjectBuild(gtkMock *gtk.Mock, objectName string, ret interface{}) *gtk.MockBuilder {
	builderMock := &gtk.MockBuilder{}

	fileContent := "this is an interface description for the object"

	gtkMock.On("BuilderNew").Return(builderMock, nil).Once()
	builderMock.On("AddFromString", fileContent).Return(nil)
	builderMock.On("GetObject", objectName).Return(ret, nil)
	return builderMock
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

	s.gtkMock.On("ApplicationNew", mock.Anything, mock.Anything).Return(appMock, nil)

	gdkMock := &gdk.Mock{}
	ka := fixedKeyAccess(
		fixedKeyEntry("/home/amnesia/.ssh/id_ed25519"),
		fixedKeyEntry("/home/amnesia/.ssh/id_rsa"),
	)

	log, _ := test.NewNullLogger()
	Start(s.gtkMock, gdkMock, log, ka)

	winMock := &gtk.MockApplicationWindow{}
	winMock.On("SetApplication", appMock).Return().Once()
	winMock.On("ShowAll").Return().Once()
	winMock.On("GetAllocatedHeight").Return(42)
	winMock.On("Resize", 1, 42).Return()

	listBox := &gtk.MockBox{}
	detailsBox := &gtk.MockBox{}
	detailsBox.On("Hide").Return().Once()

	detailsRevealer := &gtk.MockRevealer{}

	builderMock := mockObjectBuild(s.gtkMock, "MainWindow", winMock)
	builderMock.On("GetObject", "keyListBox").Return(listBox, nil).Once()
	builderMock.On("GetObject", "keyDetailsBox").Return(detailsBox, nil).Once()
	builderMock.On("GetObject", "keyDetailsRevealer").Return(detailsRevealer, nil).Once()

	box1 := s.setupBuildingOfKeyEntry("/home/amnesia/.ssh/id_ed25519")
	box1.On("Connect", "clicked", mock.Anything).Return(nil).Once()
	box2 := s.setupBuildingOfKeyEntry("/home/amnesia/.ssh/id_rsa")
	box2.On("Connect", "clicked", mock.Anything).Return(nil).Once()

	listBox.On("Add", box1).Return().Once()
	listBox.On("Add", box2).Return().Once()

	stubStyleProviders(s.gtkMock, gdkMock)
	defer setupStubbedDefinitions()()

	activateEventHandler()

	winMock.AssertExpectations(s.T())
	builderMock.AssertExpectations(s.T())
	listBox.AssertExpectations(s.T())
}
