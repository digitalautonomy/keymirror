package gui

import (
	"crypto/sha1"
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
	builderKeyDetailsBoxMock := s.setupBuildingOfObject(keyDetailsBoxMock, "KeyDetails")

	keyDetailsHolder := &gtk.MockBox{}
	keyDetailsHolder.On("Add", keyDetailsBoxMock).Return().Once()
	keyDetailsHolder.On("GetChildren").Return(nil).Once()

	pathPublicKeyPath := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "publicKeyPath").Return(pathPublicKeyPath, nil).Once()

	labelPrivateKeyPath := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "privateKeyPathLabel").Return(labelPrivateKeyPath, nil).Once()
	pathPrivateKey := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "privateKeyPath").Return(pathPrivateKey, nil).Once()

	fingerprintSha1 := &gtk.MockLabel{}
	fingerprintSha256 := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "sha1Fingerprint").Return(fingerprintSha1, nil).Once()
	builderKeyDetailsBoxMock.On("GetObject", "sha256Fingerprint").Return(fingerprintSha256, nil).Once()

	keMock := &publicKeyEntryMock{}
	keMock.On("WithDigestContent", mock.Anything).Return([]byte{0xAB, 0xCD, 0x10}).Once()
	keMock.On("WithDigestContent", mock.Anything).Return([]byte{0xCC, 0x07, 0x00, 0xFF}).Once()
	keMock.On("PublicKeyLocations").Return([]string{"/a/path/to/a/public/key"}).Once()
	keMock.On("PrivateKeyLocations").Return(nil).Once()
	keMock.On("KeyType").Return(api.PublicKeyType).Once()
	pathPublicKeyPath.On("SetLabel", "/a/path/to/a/public/key").Return().Once()
	pathPublicKeyPath.On("SetTooltipText", "/a/path/to/a/public/key").Return().Once()
	labelPrivateKeyPath.On("Hide").Return().Once()
	pathPrivateKey.On("Hide").Return().Once()

	fingerprintSha1.On("SetLabel", "AB:CD:10").Return().Once()
	fingerprintSha1.On("SetTooltipText", "AB:CD:10").Return().Once()

	fingerprintSha256.On("SetLabel", "CC:07:00:FF").Return().Once()
	fingerprintSha256.On("SetTooltipText", "CC:07:00:FF").Return().Once()

	scMock := expectClassToBeAdded(keyDetailsBoxMock, "publicKey")

	u := &ui{gtk: s.gtkMock}
	u.populateKeyDetails(keMock, keyDetailsHolder)

	keyDetailsHolder.AssertExpectations(s.T())
	keMock.AssertExpectations(s.T())
	pathPublicKeyPath.AssertExpectations(s.T())
	labelPrivateKeyPath.AssertExpectations(s.T())
	pathPrivateKey.AssertExpectations(s.T())
	fingerprintSha1.AssertExpectations(s.T())
	fingerprintSha256.AssertExpectations(s.T())
	scMock.AssertExpectations(s.T())
}

func (s *guiSuite) Test_populateKeyDetails_createsTheKeyDetailsBoxAndDisplaysThePrivateKeyPath() {
	keyDetailsBoxMock := &gtk.MockBox{}
	builderKeyDetailsBoxMock := s.setupBuildingOfObject(keyDetailsBoxMock, "KeyDetails")

	keyDetailsHolder := &gtk.MockBox{}
	keyDetailsHolder.On("Add", keyDetailsBoxMock).Return().Once()
	keyDetailsHolder.On("GetChildren").Return(nil).Once()

	pathPrivateKey := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "privateKeyPath").Return(pathPrivateKey, nil).Once()

	labelPublicKeyPath := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "publicKeyPathLabel").Return(labelPublicKeyPath, nil).Once()
	pathPublicKey := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "publicKeyPath").Return(pathPublicKey, nil).Once()

	labelSha1Fingerprint := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "sha1FingerprintLabel").Return(labelSha1Fingerprint, nil).Once()
	fingerprintSha1 := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "sha1Fingerprint").Return(fingerprintSha1, nil).Once()

	labelSha256Fingerprint := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "sha256FingerprintLabel").Return(labelSha256Fingerprint, nil).Once()
	fingerprintSha256 := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "sha256Fingerprint").Return(fingerprintSha256, nil).Once()

	keMock := &keyEntryMock{}
	keMock.On("PublicKeyLocations").Return(nil).Once()
	keMock.On("PrivateKeyLocations").Return([]string{"/a/path/to/a/private/key"}).Once()
	keMock.On("KeyType").Return(api.PrivateKeyType).Once()
	pathPrivateKey.On("SetLabel", "/a/path/to/a/private/key").Return().Once()
	pathPrivateKey.On("SetTooltipText", "/a/path/to/a/private/key").Return().Once()
	labelPublicKeyPath.On("Hide").Return().Once()
	pathPublicKey.On("Hide").Return().Once()
	labelSha1Fingerprint.On("Hide").Return().Once()
	fingerprintSha1.On("Hide").Return().Once()
	labelSha256Fingerprint.On("Hide").Return().Once()
	fingerprintSha256.On("Hide").Return().Once()

	scMock := expectClassToBeAdded(keyDetailsBoxMock, "privateKey")

	u := &ui{gtk: s.gtkMock}
	u.populateKeyDetails(keMock, keyDetailsHolder)

	keyDetailsHolder.AssertExpectations(s.T())
	keMock.AssertExpectations(s.T())
	pathPrivateKey.AssertExpectations(s.T())
	labelPublicKeyPath.AssertExpectations(s.T())
	pathPublicKey.AssertExpectations(s.T())
	labelSha1Fingerprint.AssertExpectations(s.T())
	fingerprintSha1.AssertExpectations(s.T())
	labelSha256Fingerprint.AssertExpectations(s.T())
	fingerprintSha256.AssertExpectations(s.T())
	scMock.AssertExpectations(s.T())
}

func (s *guiSuite) Test_populateKeyDetails_createsTheKeyDetailsBoxAndDisplaysBothPublicAndPrivateKeyPathIfExists() {
	keyDetailsBoxMock := &gtk.MockBox{}
	builderKeyDetailsBoxMock := s.setupBuildingOfObject(keyDetailsBoxMock, "KeyDetails")

	keyDetailsHolder := &gtk.MockBox{}
	keyDetailsHolder.On("Add", keyDetailsBoxMock).Return().Once()
	keyDetailsHolder.On("GetChildren").Return(nil).Once()

	pathPublicKey := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "publicKeyPath").Return(pathPublicKey, nil).Once()

	pathPrivateKey := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "privateKeyPath").Return(pathPrivateKey, nil).Once()

	labelFingerprintSha1 := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "sha1FingerprintLabel").Return(labelFingerprintSha1, nil).Once()
	fingerprintSha1 := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "sha1Fingerprint").Return(fingerprintSha1, nil).Once()

	labelFingerprintSha256 := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "sha256FingerprintLabel").Return(labelFingerprintSha256, nil).Once()
	fingerprintSha256 := &gtk.MockLabel{}
	builderKeyDetailsBoxMock.On("GetObject", "sha256Fingerprint").Return(fingerprintSha256, nil).Once()

	keMock := &keyEntryMock{}
	keMock.On("PublicKeyLocations").Return([]string{"/a/path/to/a/public/key"}).Once()
	keMock.On("PrivateKeyLocations").Return([]string{"/a/path/to/a/private/key"}).Once()
	keMock.On("KeyType").Return(api.PairKeyType).Once()
	pathPublicKey.On("SetLabel", "/a/path/to/a/public/key").Return().Once()
	pathPublicKey.On("SetTooltipText", "/a/path/to/a/public/key").Return().Once()
	pathPrivateKey.On("SetLabel", "/a/path/to/a/private/key").Return().Once()
	pathPrivateKey.On("SetTooltipText", "/a/path/to/a/private/key").Return().Once()
	labelFingerprintSha1.On("Hide").Return().Once()
	fingerprintSha1.On("Hide").Return().Once()
	labelFingerprintSha256.On("Hide").Return().Once()
	fingerprintSha256.On("Hide").Return().Once()

	scMock := expectClassToBeAdded(keyDetailsBoxMock, "keyPair")

	u := &ui{gtk: s.gtkMock}
	u.populateKeyDetails(keMock, keyDetailsHolder)

	keyDetailsHolder.AssertExpectations(s.T())
	keMock.AssertExpectations(s.T())
	pathPublicKey.AssertExpectations(s.T())
	pathPrivateKey.AssertExpectations(s.T())
	labelFingerprintSha1.AssertExpectations(s.T())
	fingerprintSha1.AssertExpectations(s.T())
	fingerprintSha256.AssertExpectations(s.T())
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

func (s *guiSuite) Test_formatFingerprint_returnsAnUpperCaseHexadecimalStringWithColons() {
	f := []byte{}
	expected := ""
	s.Equal(expected, formatFingerprint(f))

	f = []byte{0}
	expected = "00"
	s.Equal(expected, formatFingerprint(f))

	f = []byte{8}
	expected = "08"
	s.Equal(expected, formatFingerprint(f))

	f = []byte{0xfe}
	expected = "FE"
	s.Equal(expected, formatFingerprint(f))

	f = []byte{0, 1, 32, 0x67, 0, 7, 0xfc, 0}
	expected = "00:01:20:67:00:07:FC:00"
	s.Equal(expected, formatFingerprint(f))
}

func (s *guiSuite) Test_keyDetails_displayFingerprint_calculateTheFingerprintAndDisplaysIt() {
	keyMock := &publicKeyEntryMock{}
	var calledWithFunc *func([]byte) []byte
	keyMock.On("WithDigestContent", mock.AnythingOfType("func([]uint8) []uint8")).Return(
		[]byte("something")).Run(func(a mock.Arguments) {
		ff := a.Get(0).(func([]byte) []byte)
		calledWithFunc = &ff
	})

	builderMock := &gtk.MockBuilder{}

	kd := &keyDetails{
		builder: &builder{builderMock},
		key:     keyMock,
	}

	labelMock := &gtk.MockLabel{}
	builderMock.On("GetObject", "sha1Fingerprint").Return(labelMock, nil).Once()
	labelMock.On("SetLabel", "73:6F:6D:65:74:68:69:6E:67").Return().Once()
	labelMock.On("SetTooltipText", "73:6F:6D:65:74:68:69:6E:67").Return().Once()

	kd.displayFingerprint("a row", "sha1Fingerprint", returningSlice20(sha1.Sum))

	labelMock.AssertExpectations(s.T())
	builderMock.AssertExpectations(s.T())
	keyMock.AssertExpectations(s.T())

	s.NotNil(calledWithFunc)
	s.Equal([]byte{0x2a, 0xae, 0x6c, 0x35, 0xc9, 0x4f, 0xcf, 0xb4, 0x15, 0xdb, 0xe9, 0x5f, 0x40, 0x8b, 0x9c, 0xe9, 0x1e, 0xe8, 0x46, 0xed}, (*calledWithFunc)([]byte("hello world")))
	s.Equal([]byte{0x0, 0x78, 0xbb, 0x8e, 0x5c, 0x9d, 0x8a, 0xbf, 0x7f, 0x1e, 0x4e, 0x14, 0xc8, 0x7d, 0x90, 0x23, 0x23, 0x5b, 0x62, 0x30}, (*calledWithFunc)([]byte("goodbye world")))
}

func (s *guiSuite) Test_keyDetails_displayFingerprint_hideTheFingerprintRow_ifDisplayAPrivateKey() {
	keyMock := &keyEntryMock{}

	builderMock := &gtk.MockBuilder{}
	labelFingerprintSha1 := &gtk.MockLabel{}
	fingerprintSha1 := &gtk.MockLabel{}

	kd := &keyDetails{
		builder: &builder{builderMock},
		key:     keyMock,
	}

	labelFingerprintSha1.On("Hide").Return().Once()
	builderMock.On("GetObject", "labelFingerprintSha1").Return(labelFingerprintSha1, nil).Maybe()
	fingerprintSha1.On("Hide").Return().Once()
	builderMock.On("GetObject", "fingerprintSha1").Return(fingerprintSha1, nil).Maybe()

	kd.displayFingerprint("labelFingerprintSha1", "fingerprintSha1", func([]byte) []byte {
		return []byte{}
	})

	labelFingerprintSha1.AssertExpectations(s.T())
}
