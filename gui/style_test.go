package gui

import (
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/prashantv/gostub"
	"testing/fstest"
)

func (s *guiSuite) Test_createStyleProviderFrom_CreatesAStyleProviderFromAStyleFile() {
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

	sp := createStyleProviderFrom(gtkMock, filename)

	s.Equal(cssProviderMock, sp)
	gtkMock.AssertExpectations(s.T())
	cssProviderMock.AssertExpectations(s.T())
}
