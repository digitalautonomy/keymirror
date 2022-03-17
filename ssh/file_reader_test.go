package ssh

import (
	"os"
	"path"
	"path/filepath"
)

func (s *sshSuite) createFileWithContent(dir, fileName, content string) {
	e := os.WriteFile(path.Join(dir, fileName), []byte(content), 0666)
	s.Nil(e)
}

func (s *sshSuite) createEmptyFile(dir, fileName string) {
	s.createFileWithContent(dir, fileName, "")
}

func (s *sshSuite) Test_checkIfFileContainsAPublicRSAKey_returnsAnErrorWhenFileDoesntExist() {
	fileName := "a-file-that-doesnt-exist"

	_, err := checkIfFileContainsAPublicRSAKey(fileName)

	s.True(os.IsNotExist(err), "Function should generate an error indicating the file doesn`t exist")
}

func (s *sshSuite) Test_checkIfFileContainsAPublicRSAKey_doesNotReturnErrorWhenFileExists() {
	fileName := "a file that should exist"
	s.createEmptyFile(s.tdir, fileName)

	_, err := checkIfFileContainsAPublicRSAKey(filepath.Join(s.tdir, fileName))

	s.Nil(err)
}

func (s *sshSuite) Test_checkIfFileContainsAPublicRSAKey_ReturnsFalseWhenTheFileContentIsNotAnRSAPublicKey() {
	fileName := "a file without an RSA Public Key"
	s.createFileWithContent(s.tdir, fileName, "not an RSA Public Key")

	b, err := checkIfFileContainsAPublicRSAKey(filepath.Join(s.tdir, fileName))

	s.Nil(err)
	s.False(b)
}

func (s *sshSuite) Test_checkIfFileContainsAPublicRSAKey_ReturnsTrueWhenTheFileContentIsAnRSAPublicKey() {
	fileName := "a file with a valid RSA Public Key"
	s.createFileWithContent(s.tdir, fileName, "ssh-rsa AAAAA batman@debian")

	b, err := checkIfFileContainsAPublicRSAKey(filepath.Join(s.tdir, fileName))

	s.Nil(err)
	s.True(b)
}
