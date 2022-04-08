package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3mocks/gdk"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/prashantv/gostub"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"testing/fstest"
)

func (s *guiSuite) Test_ui_createStyleProviderFrom_CreatesAStyleProviderFromAStyleFile() {
	cssProviderMock := &gtk.MockCssProvider{}
	gtkMock := &gtk.Mock{}
	gtkMock.On("CssProviderNew").Return(cssProviderMock, nil).Once()

	filename := "a-fancy-style-sheet"
	content := "css-content"
	cssProviderMock.On("LoadFromData", content).Return(nil).Once()
	definitionsMock := fstest.MapFS{
		"definitions/styles/a-fancy-style-sheet.css": &fstest.MapFile{
			Data: []byte(content),
		},
	}
	defer gostub.StubFunc(&getDefinitions, definitionsMock).Reset()

	log, h := test.NewNullLogger()
	log.Level = logrus.TraceLevel
	testUI := &ui{
		gtk: gtkMock,
		log: log,
	}

	sp := testUI.createStyleProviderFrom(filename)

	s.Equal(cssProviderMock, sp)
	gtkMock.AssertExpectations(s.T())
	cssProviderMock.AssertExpectations(s.T())

	s.Len(h.Entries, 1)
	s.Equal("loading CSS style", h.LastEntry().Message)
	s.Equal(logrus.DebugLevel, h.LastEntry().Level)
	s.Len(h.LastEntry().Data, 1)
	s.Equal("definitions/styles/a-fancy-style-sheet.css", h.LastEntry().Data["file"])
}

func (s *guiSuite) Test_ui_applyApplicationStyle_LoadsGlobalAndColorsStylesToTheDefaultScreen() {
	defer gostub.Stub(&gtki.STYLE_PROVIDER_PRIORITY_APPLICATION, gtki.StyleProviderPriority(77)).Reset()

	cssProvider := &gtk.MockCssProvider{}

	screenMock := &gdk.MockScreen{}

	gtkMock := &gtk.Mock{}

	gtkMock.On("CssProviderNew").Return(cssProvider, nil).Twice()

	cssProvider.On("LoadFromData", "global content").Return(nil).Once()
	cssProvider.On("LoadFromData", "light colors content").Return(nil).Once()

	definitionsMock := fstest.MapFS{
		"definitions/styles/global.css": &fstest.MapFile{
			Data: []byte("global content"),
		},
		"definitions/styles/colors_light.css": &fstest.MapFile{
			Data: []byte("light colors content"),
		},
	}
	defer gostub.StubFunc(&getDefinitions, definitionsMock).Reset()

	gtkMock.On("AddProviderForScreen", screenMock, cssProvider, uint(gtki.STYLE_PROVIDER_PRIORITY_APPLICATION)).Return().Twice()

	gdkMock := &gdk.Mock{}
	gdkMock.On("ScreenGetDefault").Return(screenMock, nil).Once()

	log, _ := test.NewNullLogger()
	testUI := &ui{
		gtk: gtkMock,
		gdk: gdkMock,
		log: log,
	}

	testUI.applyApplicationStyle()

	gtkMock.AssertExpectations(s.T())
	gdkMock.AssertExpectations(s.T())
	cssProvider.AssertExpectations(s.T())
}
