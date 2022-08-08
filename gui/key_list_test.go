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

func (ke *keyEntryMock) Size() int {
	returns := ke.Called()
	return returns.Int(0)
}

func (ke *keyEntryMock) Algorithm() api.Algorithm {
	returns := ke.Called()
	return ret[api.Algorithm](returns, 0)
}

type publicKeyEntryMock struct {
	keyEntryMock
}

func (pk *publicKeyEntryMock) WithDigestContent(f func([]byte) []byte) []byte {
	returns := pk.Called(f)
	return ret[[]byte](returns, 0)
}

func (pk *publicKeyEntryMock) UserID() string {
	returns := pk.Called()
	return returns.String(0)
}

func (s *guiSuite) Test_createKeyEntryBoxFrom_CreatesAGTKIBoxWithTheGivenASSHKeyEntry() {
	box := s.setupBuildingOfKeyEntry("/home/amnesia/id_ed25519.pub", "Ed25519")

	onWindowChangeCalled := 0
	u := &ui{gtk: s.gtkMock, onWindowSizeChange: func() {
		onWindowChangeCalled++
	}}

	keyEntry := &keyEntryMock{}
	keyEntry.On("Locations").Return([]string{"/home/amnesia/id_ed25519.pub"}).Once()
	keyEntry.On("Size").Return(0).Maybe()
	keyEntry.On("Algorithm").Return(api.Ed25519).Maybe()

	var clickedHandler func() = nil
	box.On("Connect", "clicked", mock.Anything).Return(nil).Once().Run(func(a mock.Arguments) {
		clickedHandler = a.Get(1).(func())
	})

	detailsBoxMock := &gtk.MockBox{}

	detailsRevMock := &gtk.MockRevealer{}
	actualGtkBox := u.createKeyEntryBoxFrom(keyEntry, detailsBoxMock, detailsRevMock)

	keyDetailsBoxMock := &gtk.MockBox{}
	builderKeyDetailsBoxMock := s.setupBuildingOfObject(keyDetailsBoxMock, "KeyDetails")

	detailsBoxMock.On("Add", keyDetailsBoxMock).Return().Once()
	detailsBoxMock.On("GetChildren").Return(nil).Once()
	pathPublicKey := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "publicKeyPath").Return(pathPublicKey, nil).Once()

	s.addLabelsThatShouldHide(builderKeyDetailsBoxMock,
		"privateKeyPath",
		"privateKeyPathLabel",
		"passwordProtectedLabel",
		"sha1FingerprintLabel",
		"sha1Fingerprint",
		"sha256FingerprintLabel",
		"sha256Fingerprint",
		"userIDLabel",
		"userID",
	)

	keyEntry.On("PublicKeyLocations").Return([]string{"/a/path/to/a/public/key"}).Once()
	keyEntry.On("PrivateKeyLocations").Return(nil).Once()
	keyEntry.On("KeyType").Return(api.PublicKeyType).Maybe()

	pathPublicKey.On("SetLabel", "/a/path/to/a/public/key").Return().Once()
	pathPublicKey.On("SetTooltipText", "/a/path/to/a/public/key").Return().Once()

	detailsRevMock.On("Show").Return().Once()
	detailsRevMock.On("SetRevealChild", true).Return().Once()

	notificationMessage := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "notification").Return(notificationMessage, nil).Once()
	notificationMessage.On("SetLabel", "(no private key available)").Return().Once()

	properties := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "algorithm").Return(properties, nil).Once()
	properties.On("SetLabel", "Ed25519").Return().Once()

	scMock1 := expectClassToBeAdded(box, "current")
	scMock2 := expectClassToBeAdded(keyDetailsBoxMock, "publicKey")

	clickedHandler()

	s.Equal(1, onWindowChangeCalled)
	s.Equal(box, actualGtkBox)
	keyEntry.AssertExpectations(s.T())
	notificationMessage.AssertExpectations(s.T())
	detailsBoxMock.AssertExpectations(s.T())
	detailsRevMock.AssertExpectations(s.T())
	scMock1.AssertExpectations(s.T())
	scMock2.AssertExpectations(s.T())
	properties.AssertExpectations(s.T())
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

func fixedKeyEntry(location string, algo api.Algorithm) api.KeyEntry {
	ke := &keyEntryMock{}
	ke.On("Locations").Return([]string{location}).Maybe()
	ke.On("Algorithm").Return(algo).Maybe()
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

func (s *guiSuite) setupBuildingOfKeyEntry(path, algo string) *gtk.MockButton {
	label := &gtk.MockLabel{}
	label.On("SetLabel", path).Return().Once()
	algorithm := &gtk.MockLabel{}
	algorithm.On("SetLabel", algo).Return().Once()
	box := &gtk.MockButton{}
	b := s.setupBuildingOfObject(box, "KeyListEntry")
	b.On("GetObject", "keyListEntryLabel").Return(label, nil).Once()
	b.On("GetObject", "algorithmLabel").Return(algorithm, nil).Once()
	s.addObjectToAssert(box)
	return box
}

func (s *guiSuite) Test_populateListWithKeyEntries_IfThereAreKeyEntriesAddsThemIntoAGTKBoxWithoutCallingOnNoKeysFunctionPassedInParameter() {
	ka := fixedKeyAccess(
		fixedKeyEntry("/home/amnesia/.ssh/id_rsa", api.RSA),
		fixedKeyEntry("/home/amnesia/.ssh/id_ed25519", api.Ed25519),
		fixedKeyEntry("/home/amnesia/.ssh/id_dsa", api.DSA),
	)

	box1 := s.setupBuildingOfKeyEntry("/home/amnesia/.ssh/id_rsa", "RSA")
	box2 := s.setupBuildingOfKeyEntry("/home/amnesia/.ssh/id_ed25519", "Ed25519")
	box3 := s.setupBuildingOfKeyEntry("/home/amnesia/.ssh/id_dsa", "DSA")

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
