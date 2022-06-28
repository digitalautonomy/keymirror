package ssh

import (
	"fmt"
	"github.com/digitalautonomy/keymirror/api"
	"github.com/prashantv/gostub"
	"math/rand"
	"os"
	"path"
)

func (s *sshSuite) Test_access_AllKeys_ReturnsAnEmptyKeyEntryListIfCanNotFindSSHDirectory() {
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	a, _ := accessWithTestLogging()
	keys := a.AllKeys()

	s.Empty(keys)
}

func (s *sshSuite) Test_access_AllKeys_ReturnsAnEmptyKeyEntryListIfSSHDirectoryHasNoFiles() {
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	s.Nil(os.Mkdir(path.Join(s.tdir, ".ssh"), 0755))
	a, _ := accessWithTestLogging()
	keys := a.AllKeys()

	s.Empty(keys)
}

func (s *sshSuite) Test_access_AllKeys_ReturnsAnEmptyKeyEntryListIfSSHDirectoryHasNoKeyFiles() {
	sshDirectory := path.Join(s.tdir, ".ssh")
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	s.Nil(os.Mkdir(sshDirectory, 0755))
	r := rand.Int()
	files := []string{"is-not-a-key-file", fmt.Sprintf("neither-this-one-%d", r)}
	for _, f := range files {
		s.createFileWithContent(sshDirectory, f, "some content")
	}
	s.createEmptyFile(sshDirectory, "empty-file")
	a, _ := accessWithTestLogging()
	keys := a.AllKeys()

	s.Empty(keys)
}

func (s *sshSuite) Test_access_AllKeys_ReturnsAKeyEntryListOfPrivateKeysIfSSHDirectoryHasOnlyPrivateKeyFiles() {
	sshDirectory := path.Join(s.tdir, ".ssh")
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	s.Nil(os.Mkdir(sshDirectory, 0755))

	r := rand.Int()
	privateKeyFile1 := "private-key"
	privateKeyFile2 := fmt.Sprintf("is-a-private-key-%d", r)
	privateKeyFile3 := "ed25519-key"
	s.createFileWithContent(sshDirectory, privateKeyFile1, correctRSASSHPrivateKey)
	s.createFileWithContent(sshDirectory, privateKeyFile2, correctRSASSHPrivateKeyOther)
	s.createFileWithContent(sshDirectory, privateKeyFile3, correctEd25519PrivateKey)
	s.createEmptyFile(sshDirectory, "empty-file")

	a, _ := accessWithTestLogging()
	s.ElementsMatch([]api.KeyEntry{
		createPrivateKeyRepresentation(path.Join(sshDirectory, privateKeyFile1)),
		createPrivateKeyRepresentation(path.Join(sshDirectory, privateKeyFile2)),
		createPrivateKeyRepresentation(path.Join(sshDirectory, privateKeyFile3)),
	}, a.AllKeys())
}

func (s *sshSuite) Test_access_AllKeys_ReturnsAKeyEntryListOfPublicKeysIfSSHDirectoryHasOnlyPublicKeyFiles() {
	sshDirectory := path.Join(s.tdir, ".ssh")
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	s.Nil(os.Mkdir(sshDirectory, 0755))

	r := rand.Int()
	publicKeyFile1 := "ssh-rsa.pub"
	publicKeyFile2 := fmt.Sprintf("ssh-rsa-%d.pub", r)
	publicKeyFile3 := "ed25519.pub"
	s.createFileWithContent(sshDirectory, publicKeyFile1, "ssh-rsa BBBB batman@debian")
	s.createFileWithContent(sshDirectory, publicKeyFile2, "ssh-rsa AAAA robin@debian")
	s.createFileWithContent(sshDirectory, publicKeyFile3, "ssh-ed25519 CCC alfred@debian")
	s.createEmptyFile(sshDirectory, "empty-file")

	a, _ := accessWithTestLogging()
	s.ElementsMatch([]api.KeyEntry{
		createPublicKeyRepresentation(path.Join(sshDirectory, publicKeyFile1)),
		createPublicKeyRepresentation(path.Join(sshDirectory, publicKeyFile2)),
		createPublicKeyRepresentation(path.Join(sshDirectory, publicKeyFile3)),
	}, a.AllKeys())
}

func (s *sshSuite) Test_access_AllKeys_ReturnsAKeyEntryListOfKeypairsIfSSHDirectoryHasOnlyMatchingPublicAndPrivateKeys() {
	sshDirectory := path.Join(s.tdir, ".ssh")
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	s.Nil(os.Mkdir(sshDirectory, 0755))

	r := rand.Int()
	matchingPrivateKeyFile1 := "match-key"
	matchingPrivateKeyFile2 := fmt.Sprintf("is-a-match-key-%d", r)
	matchingPrivateKeyFile3 := "match-ed25519-key"
	s.createFileWithContent(sshDirectory, matchingPrivateKeyFile1, correctRSASSHPrivateKey)
	s.createFileWithContent(sshDirectory, matchingPrivateKeyFile2, correctRSASSHPrivateKeyOther)
	s.createFileWithContent(sshDirectory, matchingPrivateKeyFile3, correctEd25519PrivateKey)
	matchingPublicKeyFile1 := "match-key.pub"
	matchingPublicKeyFile2 := fmt.Sprintf("is-a-match-key-%d.pub", r)
	matchingPublicKeyFile3 := "match-ed25519-key.pub"
	s.createFileWithContent(sshDirectory, matchingPublicKeyFile1, "ssh-rsa BBBB batman@debian")
	s.createFileWithContent(sshDirectory, matchingPublicKeyFile2, "ssh-rsa AAAA robin@debian")
	s.createFileWithContent(sshDirectory, matchingPublicKeyFile3, "ssh-ed25519 CCCC alfred@debian")
	s.createEmptyFile(sshDirectory, "empty-file")

	a, _ := accessWithTestLogging()
	s.ElementsMatch([]api.KeyEntry{
		createKeypairRepresentation(createPrivateKeyRepresentation(path.Join(sshDirectory, matchingPrivateKeyFile1)), createPublicKeyRepresentation(path.Join(sshDirectory, matchingPublicKeyFile1))),
		createKeypairRepresentation(createPrivateKeyRepresentation(path.Join(sshDirectory, matchingPrivateKeyFile2)), createPublicKeyRepresentation(path.Join(sshDirectory, matchingPublicKeyFile2))),
		createKeypairRepresentation(createPrivateKeyRepresentation(path.Join(sshDirectory, matchingPrivateKeyFile3)), createPublicKeyRepresentation(path.Join(sshDirectory, matchingPublicKeyFile3))),
	}, a.AllKeys())
}

func (s *sshSuite) Test_access_AllKeys_ReturnsAKeyEntryListIfSSHDirectoryPublicAndPrivateKeys() {
	sshDirectory := path.Join(s.tdir, ".ssh")
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	s.Nil(os.Mkdir(sshDirectory, 0755))

	r := rand.Int()
	matchingPrivateKey := "match-key"
	lonelyPrivateKeyFile := fmt.Sprintf("a-private-key-%d", r)
	s.createFileWithContent(sshDirectory, matchingPrivateKey, correctRSASSHPrivateKey)
	s.createFileWithContent(sshDirectory, lonelyPrivateKeyFile, correctRSASSHPrivateKeyOther)
	matchingPublicKey := "match-key.pub"
	lonelyPublicKeyFile := fmt.Sprintf("a-public-key-%d.pub", r)
	s.createFileWithContent(sshDirectory, matchingPublicKey, "ssh-rsa BBBB batman@debian")
	s.createFileWithContent(sshDirectory, lonelyPublicKeyFile, "ssh-rsa AAAA robin@debian")
	s.createEmptyFile(sshDirectory, "empty-file")

	a, _ := accessWithTestLogging()
	s.ElementsMatch([]api.KeyEntry{
		createPrivateKeyRepresentation(path.Join(sshDirectory, lonelyPrivateKeyFile)),
		createPublicKeyRepresentation(path.Join(sshDirectory, lonelyPublicKeyFile)),
		createKeypairRepresentation(createPrivateKeyRepresentation(path.Join(sshDirectory, matchingPrivateKey)), createPublicKeyRepresentation(path.Join(sshDirectory, matchingPublicKey))),
	}, a.AllKeys())
}