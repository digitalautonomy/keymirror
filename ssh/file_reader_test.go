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

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAnEmptyListIfAnEmptyListIsProvided() {
	fileNameList := []string{}

	selected := selectFilesContainingRSAPublicKeys(fileNameList)

	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAnEmptyListIfAListWithANonExistingFileIsProvided() {
	fileNameList := []string{"File that doesn't exist"}

	selected := selectFilesContainingRSAPublicKeys(fileNameList)

	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAnEmptyListIfAListWithAnEmptyFileIsProvided() {
	// Given
	fileName := "Empty file"
	s.createEmptyFile(s.tdir, fileName)
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	selected := selectFilesContainingRSAPublicKeys(fileNameList)

	// Then
	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAnEmptyListIfAListWithAFileThatDoesntContainAnRSAPublicKeyIsProvided() {
	// Given
	fileName := "Empty file"
	s.createFileWithContent(s.tdir, fileName, "not a RSA public key")
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	selected := selectFilesContainingRSAPublicKeys(fileNameList)

	// Then
	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAListWithOneFileNameIfAListWithAFileThatContainsAnRSAPublicKeyIsProvided() {
	// Given
	fileName := "File-with-content"
	s.createFileWithContent(s.tdir, fileName, "ssh-rsa AAAAA batman@debian")
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	selected := selectFilesContainingRSAPublicKeys(fileNameList)

	// Then
	s.Equal(selected, []string{filepath.Join(s.tdir, fileName)})
}

func (s *sshSuite) withDirectory(names ...string) []string {
	return transform(names, func(name string) string {
		return filepath.Join(s.tdir, name)
	})
}

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAListWithSeveralFileNamesThatConteinsRSAKey() {
	// Given
	s.createFileWithContent(s.tdir, "key_file1", "ssh-rsa AAAA batman@debian")
	s.createFileWithContent(s.tdir, "key_file2", "ssh-ecdsa AAAA batman@debian")
	s.createFileWithContent(s.tdir, "key_file3", "ssh-rsa AAAA robin@debian")
	s.createEmptyFile(s.tdir, "key_file4")

	fileList := s.withDirectory("key_file1", "key_file2", "key_file3", "key_file4", "key_file5")

	// When
	selected := selectFilesContainingRSAPublicKeys(fileList)

	// Then
	expected := s.withDirectory("key_file1", "key_file3")
	s.Equal(expected, selected)
}
