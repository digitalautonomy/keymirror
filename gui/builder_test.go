package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/prashantv/gostub"
	"testing/fstest"
)

// what other kinds of functionality might we need from builders
// - error handling
// - we need to be able to look up other objects with other IDs, not just the main ID
// - we need to be able to connect signals at some point

func (s *guiSuite) Test_buildObjectFrom_CreatesAGTKObjectFromAnInterfaceDefinitionFile() {
	object := "MainWindow"
	gtkiMock := &gtk.Mock{}
	u := &ui{
		gtk: gtkiMock,
	}
	builderMock := &gtk.MockBuilder{}
	mainWindowMock := &gtk.MockApplicationWindow{}

	fileContent := "this is an interface description for the main window"

	definitionsMock := fstest.MapFS{
		"definitions/interface/MainWindow.xml": &fstest.MapFile{
			Data: []byte(fileContent),
		},
	}
	defer gostub.StubFunc(&getDefinitions, definitionsMock).Reset()

	gtkiMock.On("BuilderNew").Return(builderMock, nil)
	builderMock.On("AddFromString", fileContent).Return(nil)
	builderMock.On("GetObject", "MainWindow").Return(mainWindowMock, nil)

	w := buildObjectFrom[gtki.ApplicationWindow](u, object)

	s.Equal(mainWindowMock, w)

	gtkiMock.AssertExpectations(s.T())
	builderMock.AssertExpectations(s.T())
}

func (s *guiSuite) Test_buildObjectFrom_CreatesAnotherGTKObjectFromAnInterfaceDefinitionFile() {
	object := "Dialog"
	gtkiMock := &gtk.Mock{}
	u := &ui{
		gtk: gtkiMock,
	}
	builderMock := &gtk.MockBuilder{}
	dialogMock := &gtk.MockDialog{}

	fileContent := "content from other file"

	definitionsMock := fstest.MapFS{
		"definitions/interface/Dialog.xml": &fstest.MapFile{
			Data: []byte(fileContent),
		},
	}
	defer gostub.StubFunc(&getDefinitions, definitionsMock).Reset()

	gtkiMock.On("BuilderNew").Return(builderMock, nil)
	builderMock.On("AddFromString", fileContent).Return(nil)
	builderMock.On("GetObject", object).Return(dialogMock, nil)

	w := buildObjectFrom[gtki.Dialog](u, object)

	s.Equal(dialogMock, w)

	gtkiMock.AssertExpectations(s.T())
	builderMock.AssertExpectations(s.T())
}

func (s *guiSuite) Test_getDefinitions_returnsTheInternalDefinitionsByDefault() {
	s.Equal(definitionsInternal, getDefinitions())
}
