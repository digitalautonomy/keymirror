package ssh

import "path/filepath"

func (s *sshSuite) Test_filesContainingEd25519PrivateKeys_ReturnsAnEmptyListIfAnEmptyListIsProvided() {
	fileNameList := []string{}

	a, _ := accessWithTestLogging()
	selected := a.filesContainingEd25519PrivateKeys(fileNameList)

	s.Empty(selected)
}

func (s *sshSuite) Test_filesContainingEd25519PrivateKeys_ReturnsAnEmptyListIfAListWithANonExistingFileIsProvided() {
	fileNameList := []string{"File that doesn't exist"}

	a, _ := accessWithTestLogging()
	selected := a.filesContainingEd25519PrivateKeys(fileNameList)

	s.Empty(selected)
}

func (s *sshSuite) Test_filesContainingEd25519PrivateKeys_ReturnsAnEmptyListIfAListWithAnEmptyFileIsProvided() {
	// Given
	fileName := "Empty file"
	s.createEmptyFile(s.tdir, fileName)
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	a, _ := accessWithTestLogging()
	selected := a.filesContainingEd25519PrivateKeys(fileNameList)

	// Then
	s.Empty(selected)
}

func (s *sshSuite) Test_filesContainingEd25519PrivateKeys_ReturnsAnEmptyListIfAListWithAFileThatDoesntContainAnEd25519PrivateKeyIsProvided() {
	// Given
	fileName := "Empty file"
	s.createFileWithContent(s.tdir, fileName, "not a Ed25519 private key")
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	a, _ := accessWithTestLogging()
	selected := a.filesContainingEd25519PrivateKeys(fileNameList)

	// Then
	s.Empty(selected)
}

func (s *sshSuite) Test_filesContainingEd25519PrivateKeys_ReturnsAListWithOneFileNameIfAListWithAFileThatContainsAnEd25519PrivateKeyIsProvided() {
	// Given
	fileName := "File-with-content"
	s.createFileWithContent(s.tdir, fileName, correctEd25519PrivateKey)
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	a, _ := accessWithTestLogging()
	selected := a.filesContainingEd25519PrivateKeys(fileNameList)

	// Then
	s.Equal(selected, []string{filepath.Join(s.tdir, fileName)})
}

const correctEd25519PrivateKeyOther = `
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACCRcAFuKgCAnlEuMGswxk18tn2JVXH+7OkVSDbBC2WQ2gAAAJCs0I2+rNCN
vgAAAAtzc2gtZWQyNTUxOQAAACCRcAFuKgCAnlEuMGswxk18tn2JVXH+7OkVSDbBC2WQ2g
AAAED9p1K4JP1ykaLj705pfax2AVvTXryKkJxEXkp3eIuLm5FwAW4qAICeUS4wazDGTXy2
fYlVcf7s6RVINsELZZDaAAAACmZhdXN0b0BDQUQBAgM=
-----END OPENSSH PRIVATE KEY-----
`

func (s *sshSuite) Test_filesContainingEd25519PrivateKeys_ReturnsAListWithSeveralFileNamesThatContainsEd25519Key() {
	// Given
	s.createFileWithContent(s.tdir, "key_file1", correctECDSASSHPrivateKey)
	s.createFileWithContent(s.tdir, "key_file2", correctEd25519PrivateKey)
	s.createEmptyFile(s.tdir, "key_file3")
	s.createFileWithContent(s.tdir, "key_file4", correctEd25519PrivateKeyOther)

	fileList := s.withDirectory("key_file1", "key_file2", "key_file3", "key_file4", "key_file5")

	// When
	a, _ := accessWithTestLogging()
	selected := a.filesContainingEd25519PrivateKeys(fileList)

	// Then
	expected := s.withDirectory("key_file2", "key_file4")
	s.Equal(expected, selected)
}
