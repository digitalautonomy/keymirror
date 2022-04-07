package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3mocks/gdk"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/prashantv/gostub"
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

	testUI := &ui{
		gtk: gtkMock,
	}
	sp := testUI.createStyleProviderFrom(filename)

	s.Equal(cssProviderMock, sp)
	gtkMock.AssertExpectations(s.T())
	cssProviderMock.AssertExpectations(s.T())
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

	testUI := &ui{
		gtk: gtkMock,
		gdk: gdkMock,
	}

	testUI.applyApplicationStyle()

	gtkMock.AssertExpectations(s.T())
	gdkMock.AssertExpectations(s.T())
	cssProvider.AssertExpectations(s.T())
}
