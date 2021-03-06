package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/digitalautonomy/keymirror/api"
	"github.com/digitalautonomy/keymirror/i18n"
	"github.com/stretchr/testify/mock"
)

type getter interface {
	Get(int) interface{}
}

func safeCast[T any](v interface{}) T {
	ret, _ := v.(T)
	return ret
}

func ret[T any](args getter, index int) T {
	return safeCast[T](args.Get(index))
}

type keyEntryMock struct {
	mock.Mock
}

func (ke *keyEntryMock) Locations() []string {
	returns := ke.Called()
	return ret[[]string](returns, 0)
}

func (ke *keyEntryMock) PublicKeyLocations() []string {
	returns := ke.Called()
	return ret[[]string](returns, 0)
}

func (ke *keyEntryMock) PrivateKeyLocations() []string {
	returns := ke.Called()
	return ret[[]string](returns, 0)
}

func (ke *keyEntryMock) KeyType() api.KeyType {
	returns := ke.Called()
	return ret[api.KeyType](returns, 0)
}

func (s *guiSuite) Test_createKeyEntryBoxFrom_CreatesAGTKIBoxWithTheGivenASSHKeyEntry() {
	box := s.setupBuildingOfKeyEntry("/home/amnesia/id_rsa.pub")

	onWindowChangeCalled := 0
	u := &ui{gtk: s.gtkMock, onWindowSizeChange: func() {
		onWindowChangeCalled++
	}}

	keyEntry := &keyEntryMock{}
	keyEntry.On("Locations").Return([]string{"/home/amnesia/id_rsa.pub"}).Once()

	var clickedHandler func() = nil
	box.On("Connect", "clicked", mock.Anything).Return(nil).Once().Run(func(a mock.Arguments) {
		clickedHandler = a.Get(1).(func())
	})

	detailsBoxMock := &gtk.MockBox{}

	detailsRevMock := &gtk.MockRevealer{}
	actualGtkBox := u.createKeyEntryBoxFrom(keyEntry, detailsBoxMock, detailsRevMock)

	s.Equal(box, actualGtkBox)

	keyEntry.AssertExpectations(s.T())

	keyDetailsBoxMock := &gtk.MockBox{}
	builder := s.setupBuildingOfObject(keyDetailsBoxMock, "KeyDetails")
	detailsBoxMock.On("Add", keyDetailsBoxMock).Return().Once()
	detailsBoxMock.On("GetChildren").Return(nil).Once()
	publicKeyPathLabel := &gtk.MockLabel{}
	builder.On("GetObject", "publicKeyPath").Return(publicKeyPathLabel, nil).Once()
	privateKeyRow := &gtk.MockBox{}
	builder.On("GetObject", "keyDetailsPrivateKeyRow").Return(privateKeyRow, nil).Once()
	keyEntry.On("PublicKeyLocations").Return([]string{"/a/path/to/a/public/key"}).Once()
	keyEntry.On("PrivateKeyLocations").Return(nil).Once()
	keyEntry.On("KeyType").Return(api.PublicKeyType).Once()
	publicKeyPathLabel.On("SetLabel", "/a/path/to/a/public/key").Return().Once()
	publicKeyPathLabel.On("SetTooltipText", "/a/path/to/a/public/key").Return().Once()
	detailsRevMock.On("Show").Return().Once()
	detailsRevMock.On("SetRevealChild", true).Return().Once()
	privateKeyRow.On("Hide").Return().Once()

	scMock1 := expectClassToBeAdded(box, "current")
	scMock2 := expectClassToBeAdded(keyDetailsBoxMock, "publicKey")

	clickedHandler()

	s.Equal(1, onWindowChangeCalled)

	keyEntry.AssertExpectations(s.T())
	detailsBoxMock.AssertExpectations(s.T())
	detailsRevMock.AssertExpectations(s.T())
	scMock1.AssertExpectations(s.T())
	scMock2.AssertExpectations(s.T())
	privateKeyRow.AssertExpectations(s.T())
}

type keyAccessMock struct {
	mock.Mock
}

func (ka *keyAccessMock) AllKeys() []api.KeyEntry {
	return ka.Called().Get(0).([]api.KeyEntry)
}

func fixedKeyAccess(keys ...api.KeyEntry) api.KeyAccess {
	ka := &keyAccessMock{}
	ka.On("AllKeys").Return(keys).Maybe()
	return ka
}

func fixedKeyEntry(location string) api.KeyEntry {
	ke := &keyEntryMock{}
	ke.On("Locations").Return([]string{location}).Maybe()
	return ke
}

func fixedPublicKeyEntry(locations ...string) api.KeyEntry {
	ke := &keyEntryMock{}
	ke.On("KeyType").Return(api.PublicKeyType).Maybe()
	ke.On("PublicKeyLocations").Return(locations).Maybe()
	return ke
}

func (s *guiSuite) setupBuildingOfObject(val interface{}, name string) *gtk.MockBuilder {
	builder := &gtk.MockBuilder{}
	s.gtkMock.On("BuilderNew").Return(builder, nil).Once()
	builder.On("AddFromString", mock.Anything).Return(nil).Once()
	builder.On("GetObject", name).Return(val, nil).Once()
	s.addObjectToAssert(builder)
	return builder
}

func (s *guiSuite) setupBuildingOfKeyEntry(path string) *gtk.MockButton {
	label := &gtk.MockLabel{}
	label.On("SetLabel", path).Return().Once()
	box := &gtk.MockButton{}
	b := s.setupBuildingOfObject(box, "KeyListEntry")
	b.On("GetObject", "keyListEntryLabel").Return(label, nil).Once()
	s.addObjectToAssert(box)
	return box
}

func (s *guiSuite) Test_populateListWithKeyEntries_IfThereAreKeyEntriesAddsThemIntoAGTKBoxWithoutCallingOnNoKeysFunctionPassedInParameter() {
	ka := fixedKeyAccess(
		fixedKeyEntry("/home/amnesia/.ssh/id_rsa"),
		fixedKeyEntry("/home/amnesia/.ssh/id_ed25519"),
		fixedKeyEntry("/home/amnesia/.ssh/id_dsa"),
	)

	box1 := s.setupBuildingOfKeyEntry("/home/amnesia/.ssh/id_rsa")
	box2 := s.setupBuildingOfKeyEntry("/home/amnesia/.ssh/id_ed25519")
	box3 := s.setupBuildingOfKeyEntry("/home/amnesia/.ssh/id_dsa")

	box1.On("Connect", "clicked", mock.Anything).Return(nil).Once()
	box2.On("Connect", "clicked", mock.Anything).Return(nil).Once()
	box3.On("Connect", "clicked", mock.Anything).Return(nil).Once()

	box := &gtk.MockBox{}
	box.On("Add", box1).Return().Once()
	box.On("Add", box2).Return().Once()
	box.On("Add", box3).Return().Once()

	u := &ui{gtk: s.gtkMock}

	called := false
	onNoKeys := func(box gtki.Box) { called = true }

	u.populateListWithKeyEntries(ka, box, nil, nil, onNoKeys)

	box.AssertExpectations(s.T())
	s.False(called)
}

func (s *guiSuite) Test_populateListWithKeyEntries_IfThereAreNoKeyEntriesExecuteOnNoKeysFunctionPassedInParameter() {
	ka := fixedKeyAccess()

	box := &gtk.MockBox{}

	u := &ui{gtk: s.gtkMock}

	called := false
	onNoKeys := func(box gtki.Box) { called = true }

	u.populateListWithKeyEntries(ka, box, nil, nil, onNoKeys)

	box.AssertExpectations(s.T())
	s.True(called)
}

func (s *guiSuite) Test_showNoAvailableKeysMessage_AddsAMessageIntoAGTKBoxWhenThereAreNoAvailableKeys() {

	sc := &gtk.MockStyleContext{}
	sc.On("AddClass", "infoMessage").Return().Once()

	label := &gtk.MockLabel{}
	label.On("GetStyleContext").Return(sc, nil).Once()

	s.gtkMock.On("LabelNew", i18n.Local("\u26A0 No keys available \u26A0")).Return(label, nil).Once()
	u := &ui{gtk: s.gtkMock}

	box := &gtk.MockBox{}
	box.On("Add", label).Return().Once()

	u.showNoAvailableKeysMessage(box)

	box.AssertExpectations(s.T())
	label.AssertExpectations(s.T())
	sc.AssertExpectations(s.T())
}
