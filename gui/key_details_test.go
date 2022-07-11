package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/digitalautonomy/keymirror/api"
	"github.com/stretchr/testify/mock"
)

type mockable interface {
	On(methodName string, arguments ...interface{}) *mock.Call
}

func expectClassToBeAdded(m mockable, class string) *gtk.MockStyleContext {
	scMock := &gtk.MockStyleContext{}
	m.On("GetStyleContext").Return(scMock, nil).Once()
	scMock.On("AddClass", class).Return().Once()
	return scMock
}

func (s *guiSuite) Test_populateKeyDetails_createsTheKeyDetailsBoxAndDisplaysThePublicKeyPath() {
	keyDetailsBoxMock := &gtk.MockBox{}
	builder := s.setupBuildingOfObject(keyDetailsBoxMock, "KeyDetails")

	keyDetailsHolder := &gtk.MockBox{}
	keyDetailsHolder.On("Add", keyDetailsBoxMock).Return().Once()
	keyDetailsHolder.On("GetChildren").Return(nil).Once()

	publicKeyPathLabel := &gtk.MockLabel{}
	builder.On("GetObject", "publicKeyPath").Return(publicKeyPathLabel, nil).Once()

	privateKeyRow := &gtk.MockBox{}
	builder.On("GetObject", "keyDetailsPrivateKeyRow").Return(privateKeyRow, nil).Once()

	keMock := &keyEntryMock{}
	keMock.On("PublicKeyLocations").Return([]string{"/a/path/to/a/public/key"}).Once()
	keMock.On("PrivateKeyLocations").Return(nil).Once()
	keMock.On("KeyType").Return(api.PublicKeyType).Once()
	publicKeyPathLabel.On("SetLabel", "/a/path/to/a/public/key").Return().Once()
	publicKeyPathLabel.On("SetTooltipText", "/a/path/to/a/public/key").Return().Once()
	privateKeyRow.On("Hide").Return().Once()

	scMock := expectClassToBeAdded(keyDetailsBoxMock, "publicKey")

	u := &ui{gtk: s.gtkMock}
	u.populateKeyDetails(keMock, keyDetailsHolder)

	keyDetailsHolder.AssertExpectations(s.T())
	keMock.AssertExpectations(s.T())
	publicKeyPathLabel.AssertExpectations(s.T())
	privateKeyRow.AssertExpectations(s.T())
	scMock.AssertExpectations(s.T())
}

func (s *guiSuite) Test_populateKeyDetails_createsTheKeyDetailsBoxAndDisplaysThePrivateKeyPath() {
	keyDetailsBoxMock := &gtk.MockBox{}
	builder := s.setupBuildingOfObject(keyDetailsBoxMock, "KeyDetails")

	keyDetailsHolder := &gtk.MockBox{}
	keyDetailsHolder.On("Add", keyDetailsBoxMock).Return().Once()
	keyDetailsHolder.On("GetChildren").Return(nil).Once()

	privateKeyPathLabel := &gtk.MockLabel{}
	builder.On("GetObject", "privateKeyPath").Return(privateKeyPathLabel, nil).Once()

	publicKeyRow := &gtk.MockBox{}
	builder.On("GetObject", "keyDetailsPublicKeyRow").Return(publicKeyRow, nil).Once()

	keMock := &keyEntryMock{}
	keMock.On("PublicKeyLocations").Return(nil).Once()
	keMock.On("PrivateKeyLocations").Return([]string{"/a/path/to/a/private/key"}).Once()
	keMock.On("KeyType").Return(api.PrivateKeyType).Once()
	privateKeyPathLabel.On("SetLabel", "/a/path/to/a/private/key").Return().Once()
	privateKeyPathLabel.On("SetTooltipText", "/a/path/to/a/private/key").Return().Once()
	publicKeyRow.On("Hide").Return().Once()

	scMock := expectClassToBeAdded(keyDetailsBoxMock, "privateKey")

	u := &ui{gtk: s.gtkMock}
	u.populateKeyDetails(keMock, keyDetailsHolder)

	keyDetailsHolder.AssertExpectations(s.T())
	keMock.AssertExpectations(s.T())
	privateKeyPathLabel.AssertExpectations(s.T())
	publicKeyRow.AssertExpectations(s.T())
	scMock.AssertExpectations(s.T())
}

func (s *guiSuite) Test_populateKeyDetails_createsTheKeyDetailsBoxAndDisplaysBothPublicAndPrivateKeyPathIfExists() {
	keyDetailsBoxMock := &gtk.MockBox{}
	builder := s.setupBuildingOfObject(keyDetailsBoxMock, "KeyDetails")

	keyDetailsHolder := &gtk.MockBox{}
	keyDetailsHolder.On("Add", keyDetailsBoxMock).Return().Once()
	keyDetailsHolder.On("GetChildren").Return(nil).Once()

	publicKeyPathLabel := &gtk.MockLabel{}
	builder.On("GetObject", "publicKeyPath").Return(publicKeyPathLabel, nil).Once()

	privateKeyPathLabel := &gtk.MockLabel{}
	builder.On("GetObject", "privateKeyPath").Return(privateKeyPathLabel, nil).Once()

	keMock := &keyEntryMock{}
	keMock.On("PublicKeyLocations").Return([]string{"/a/path/to/a/public/key"}).Once()
	keMock.On("PrivateKeyLocations").Return([]string{"/a/path/to/a/private/key"}).Once()
	keMock.On("KeyType").Return(api.PairKeyType).Once()
	publicKeyPathLabel.On("SetLabel", "/a/path/to/a/public/key").Return().Once()
	publicKeyPathLabel.On("SetTooltipText", "/a/path/to/a/public/key").Return().Once()
	privateKeyPathLabel.On("SetLabel", "/a/path/to/a/private/key").Return().Once()
	privateKeyPathLabel.On("SetTooltipText", "/a/path/to/a/private/key").Return().Once()

	scMock := expectClassToBeAdded(keyDetailsBoxMock, "keyPair")

	u := &ui{gtk: s.gtkMock}
	u.populateKeyDetails(keMock, keyDetailsHolder)

	keyDetailsHolder.AssertExpectations(s.T())
	keMock.AssertExpectations(s.T())
	publicKeyPathLabel.AssertExpectations(s.T())
	privateKeyPathLabel.AssertExpectations(s.T())
	scMock.AssertExpectations(s.T())
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
