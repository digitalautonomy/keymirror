package gui

import (
	"fmt"
	"github.com/coyim/gotk3mocks/gio"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"math/rand"
	"path"
	"testing/fstest"
)

// func makeCustomIconsAvailable()
// - it should load the resource from a path on disk
// - it should register that resource globally
// - it should add our resource prefix as a resource path for the global icon theme
// - it should take the definitions of icons from memory and unpack it
//   to a temporary file

func (s *guiSuite) Test_makeCustomIconsAvailable_loadsAResourceFromAPath() {
	file := path.Join(s.T().TempDir(), "icons.gresource")
	resource := &gio.MockResource{}
	gioMock := &gio.Mock{}
	gtkMock := &gtk.Mock{}
	gioMock.On("LoadResource", file).Return(resource, nil).Once()
	u := &ui{
		gtk: gtkMock,
		gio: gioMock,
	}

	gioMock.On("RegisterResource", resource).Return().Once()

	iconThemeMock := &gtk.MockIconTheme{}
	iconThemeMock.On("AddResourcePath", "/digital/autonomia/KeyMirror").Return().Once()

	gtkMock.On("IconThemeGetDefault").Return(iconThemeMock).Once()

	u.makeCustomIconsAvailable(file)

	gioMock.AssertExpectations(s.T())
	iconThemeMock.AssertExpectations(s.T())
}

func (s *guiSuite) Test_loadResourceDefinitions_addsTheResourcePathToTheDefaultIconTheme() {
	randomizedContent := fmt.Sprintf("this is our expected icons content %d", rand.Int())

	definitionsMock := fstest.MapFS{
		"definitions/resources/icons.gresource": &fstest.MapFile{
			Data: []byte(randomizedContent),
		},
	}
	defer gostub.StubFunc(&getDefinitions, definitionsMock).Reset()

	gioMock := &gio.Mock{}
	gtkMock := &gtk.Mock{}
	var calledWithFilePath *string
	var contentOfFile *[]byte
	var errorReadingFile *error
	gioMock.On("LoadResource", mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		p := args.String(0)
		calledWithFilePath = &p
		content, e := ioutil.ReadFile(p)
		contentOfFile = &content
		errorReadingFile = &e
	})
	gioMock.On("RegisterResource", mock.Anything).Return()
	iconThemeMock := &gtk.MockIconTheme{}
	iconThemeMock.On("AddResourcePath", mock.Anything).Return()
	gtkMock.On("IconThemeGetDefault").Return(iconThemeMock)

	u := &ui{
		gtk: gtkMock,
		gio: gioMock,
	}

	u.loadResourceDefinitions()

	gioMock.AssertExpectations(s.T())

	s.NotNil(calledWithFilePath)
	s.NotNil(contentOfFile)
	s.NotNil(errorReadingFile)

	s.NoError(*errorReadingFile)
	s.Equal(randomizedContent, string(*contentOfFile))

	s.NoFileExists(*calledWithFilePath)
}
