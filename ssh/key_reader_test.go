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
	s.False(ok)

	k = "ssh-rsa bla batman@debian"
	pub, ok := parsePublicKey(k)
	s.True(ok)
	s.Equal("ssh-rsa", pub.t)
}
