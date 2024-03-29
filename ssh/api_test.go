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
	p := a.privateKeyRepresentationsFrom([]string{
		path.Join(sshDirectory, privateKeyFile1),
		path.Join(sshDirectory, privateKeyFile2),
		path.Join(sshDirectory, privateKeyFile3),
	})
	s.ElementsMatch(p, a.AllKeys())
}

func (s *sshSuite) Test_access_AllKeys_ReturnsAKeyEntryListOfPublicKeysIfSSHDirectoryHasOnlyPublicKeyFiles() {
	sshDirectory := path.Join(s.tdir, ".ssh")
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	s.Nil(os.Mkdir(sshDirectory, 0755))

	r := rand.Int()
	publicKeyFile1 := "ssh-rsa.pub"
	publicKeyFile2 := fmt.Sprintf("ssh-rsa-%d.pub", r)
	publicKeyFile3 := "ed25519.pub"
	publicKeyFile4 := "other_ed25519.pub"
	s.createFileWithContent(sshDirectory, publicKeyFile1, "ssh-rsa BBBB batman@debian")
	s.createFileWithContent(sshDirectory, publicKeyFile2, "ssh-rsa AAAA robin@debian")
	s.createFileWithContent(sshDirectory, publicKeyFile3, "ssh-ed25519 CCCC alfred@debian")
	s.createFileWithContent(sshDirectory, publicKeyFile4, "ssh-ed25519 DDD penguin@debian")
	s.createEmptyFile(sshDirectory, "empty-file")

	a, _ := accessWithTestLogging()
	p := publicKeyRepresentationsFrom([]string{
		path.Join(sshDirectory, publicKeyFile1),
		path.Join(sshDirectory, publicKeyFile2),
		path.Join(sshDirectory, publicKeyFile3),
	})
	s.ElementsMatch([]api.KeyEntry{p[0], p[1], p[2]}, a.AllKeys())
}

func createPublicKeyRepresentationForTest(path, key string) *publicKeyRepresentation {
	return createPublicKeyRepresentationFromPublicKey(&publicKey{
		location: path,
		key:      decode(key),
	})
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
	privateKeys := a.privateKeyRepresentationsFrom([]string{
		path.Join(sshDirectory, matchingPrivateKeyFile1),
		path.Join(sshDirectory, matchingPrivateKeyFile2),
		path.Join(sshDirectory, matchingPrivateKeyFile3),
	})

	publicKeys := publicKeyRepresentationsFrom([]string{
		path.Join(sshDirectory, matchingPublicKeyFile1),
		path.Join(sshDirectory, matchingPublicKeyFile2),
		path.Join(sshDirectory, matchingPublicKeyFile3),
	})
	s.ElementsMatch([]api.KeyEntry{
		createKeypairRepresentation(privateKeys[0], publicKeys[0]),
		createKeypairRepresentation(privateKeys[1], publicKeys[1]),
		createKeypairRepresentation(privateKeys[2], publicKeys[2]),
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
	privateKeys := a.privateKeyRepresentationsFrom([]string{
		path.Join(sshDirectory, lonelyPrivateKeyFile),
		path.Join(sshDirectory, matchingPrivateKey),
	})

	publicKeys := publicKeyRepresentationsFrom([]string{
		path.Join(sshDirectory, lonelyPublicKeyFile),
		path.Join(sshDirectory, matchingPublicKey),
	})
	s.ElementsMatch([]api.KeyEntry{
		privateKeys[0],
		publicKeys[0],
		createKeypairRepresentation(privateKeys[1], publicKeys[1]),
	}, a.AllKeys())
}
