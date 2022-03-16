package ssh

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type sshSuite struct {
	suite.Suite
}

func TestSSHSuite(t *testing.T) {
	suite.Run(t, new(sshSuite))
}

func (s *sshSuite) Test_ListsAllTheFilesInSpecifiedDirectory() {
	dir := s.T().TempDir()

	r := rand.Int()

	expected := []string{"id_rsa.pub", fmt.Sprintf("id_rsa%d", r)}

	for _, f := range expected {
		file := filepath.Join(dir, f)
		err := os.WriteFile(file, []byte("some content"), 0666)

		s.Nil(err)
	}

	s.Equal(expected, listFilesIn(dir))
}

func (s *sshSuite) Test_ListsNoFilesInADirectoryThatDoesntExist() {
	s.Equal([]string{}, listFilesIn("directory that hopefully doesnt exist"))
}

func (s *sshSuite) Test_ParseAStringAsAnSSHPublicKeyRepresentation() {
	k := ""
	_, ok := parsePublicKey(k)
	s.Require().False(ok, "An empty string is not a valid SSH public key representation")

	k = "ssh-rsa bla batman@debian"
	pub, ok := parsePublicKey(k)
	s.Require().True(ok, "Should parse a valid SSH RSA public key representation")
	s.Equal("ssh-rsa", pub.algorithm)

	k = "ssh-ecdsa bla2 robin@debian"
	pub, ok = parsePublicKey(k)
	s.Require().True(ok, "Should parse a valid SSH public key representation with a different key type")
	s.Equal("ssh-ecdsa", pub.algorithm)

	k = "ssh-rsa"
	_, ok = parsePublicKey(k)
	s.Require().False(ok, "A string with only one field is not a valid SSH public key representation")

	k = "ssh-rsa  "
	_, ok = parsePublicKey(k)
	s.False(ok, "Since more than one whitespace character serve as one single separator, this example only has one column, and is thus not valid")

	k = "ssh-rsa  AAAAA foo@debian"
	_, ok = parsePublicKey(k)
	s.True(ok, "More than one whitespace character serves as one single separator between columns")

	k = "ssh-rsa\tAAAAA foo@debian"
	_, ok = parsePublicKey(k)
	s.True(ok, "A tab can be a separator for columns")

	k = "ssh-rsa   \t  \t  \t AAAAA foo@debian"
	pub, ok = parsePublicKey(k)
	s.True(ok, "A mix of tabs and spaces serve as one single separator")
	s.Equal("AAAAA", pub.key)

	k = "ssh-rsa AAQQ foo@debian foo2@debian"
	pub, ok = parsePublicKey(k)
	s.True(ok, "More than one comment is acceptable in an SSH public key")
	s.Equal("foo@debian foo2@debian", pub.comment)

	k = "ssh-rsa AAQQ"
	_, ok = parsePublicKey(k)
	s.True(ok, "An SSH public key without a comment is still acceptable")
}

func (s *sshSuite) Test_CheckIfTheTypeIdentifierIsSSHRSA() {
	pub := publicKey{}
	s.False(pub.isRSA(), "An empty key is not an RSA key")

	pub = publicKey{algorithm: rsaAlgorithm}
	s.True(pub.isRSA(), "A key with the algorithm identifier ssh-rsa is an RSA key")

	pub = publicKey{algorithm: "ssh-ecdsa"}
	s.False(pub.isRSA(), "A key with the algorithm identifier ssh-ecdsa is not an RSA key")
}
