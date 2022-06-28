package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3mocks/gtk"
)

func (s *guiSuite) Test_populateKeyDetails_createsTheKeyDetailsBoxAndDisplaysThePublicKeyPath() {
	keyDetailsBoxMock := &gtk.MockBox{}
	builder := s.setupBuildingOfObject(keyDetailsBoxMock, "KeyDetails")

	keyDetailsHolder := &gtk.MockBox{}
	keyDetailsHolder.On("Add", keyDetailsBoxMock).Return().Once()
	keyDetailsHolder.On("GetChildren").Return(nil).Once()

	publicKeyPathLabel := &gtk.MockLabel{}
	builder.On("GetObject", "publicKeyPath").Return(publicKeyPathLabel, nil).Once()

	keMock := &keyEntryMock{}
	keMock.On("PublicKeyLocations").Return([]string{"/a/path/to/a/public/key"}).Once()

	publicKeyPathLabel.On("SetLabel", "/a/path/to/a/public/key").Return().Once()

	u := &ui{gtk: s.gtkMock}
	u.populateKeyDetails(keMock, keyDetailsHolder)

	keyDetailsHolder.AssertExpectations(s.T())
	keMock.AssertExpectations(s.T())
	publicKeyPathLabel.AssertExpectations(s.T())
}

func (s *guiSuite) Test_clearAllChildrenOf_removeEachOneOfTheChildrenOfTheBox() {
	boxMock := &gtk.MockBox{}
	child1 := &gtk.MockWidget{}
	child2 := &gtk.MockWidget{}

	boxMock.On("GetChildren").Return([]gtki.Widget{child1, child2}).Once()
	boxMock.On("Remove", child1).Return().Once()
	boxMock.On("Remove", child2).Return().Once()

	clearAllChildrenOf[gtki.Widget](boxMock)

	boxMock.AssertExpectations(s.T())
}
