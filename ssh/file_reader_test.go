package ssh

import (
	"fmt"
	"github.com/digitalautonomy/keymirror/api"
	"github.com/sirupsen/logrus/hooks/test"
	"math/rand"
	"os"
	"path"
	"path/filepath"

	"github.com/prashantv/gostub"
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

	selected := filesContainingRSAPublicKeys(fileNameList)

	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAnEmptyListIfAListWithANonExistingFileIsProvided() {
	fileNameList := []string{"File that doesn't exist"}

	selected := filesContainingRSAPublicKeys(fileNameList)

	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAnEmptyListIfAListWithAnEmptyFileIsProvided() {
	// Given
	fileName := "Empty file"
	s.createEmptyFile(s.tdir, fileName)
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	selected := filesContainingRSAPublicKeys(fileNameList)

	// Then
	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAnEmptyListIfAListWithAFileThatDoesntContainAnRSAPublicKeyIsProvided() {
	// Given
	fileName := "Empty file"
	s.createFileWithContent(s.tdir, fileName, "not a RSA public key")
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	selected := filesContainingRSAPublicKeys(fileNameList)

	// Then
	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAListWithOneFileNameIfAListWithAFileThatContainsAnRSAPublicKeyIsProvided() {
	// Given
	fileName := "File-with-content"
	s.createFileWithContent(s.tdir, fileName, "ssh-rsa AAAAA batman@debian")
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	selected := filesContainingRSAPublicKeys(fileNameList)

	// Then
	s.Equal(selected, []string{filepath.Join(s.tdir, fileName)})
}

func (s *sshSuite) withDirectory(names ...string) []string {
	return transform(names, func(name string) string {
		return filepath.Join(s.tdir, name)
	})
}

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAListWithSeveralFileNamesThatContainsRSAKey() {
	// Given
	s.createFileWithContent(s.tdir, "key_file1", "ssh-rsa AAAA batman@debian")
	s.createFileWithContent(s.tdir, "key_file2", "ssh-ecdsa AAAA batman@debian")
	s.createFileWithContent(s.tdir, "key_file3", "ssh-rsa AAAA robin@debian")
	s.createEmptyFile(s.tdir, "key_file4")

	fileList := s.withDirectory("key_file1", "key_file2", "key_file3", "key_file4", "key_file5")

	// When
	selected := filesContainingRSAPublicKeys(fileList)

	// Then
	expected := s.withDirectory("key_file1", "key_file3")
	s.Equal(expected, selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPrivateKeys_ReturnsAnEmptyListIfAnEmptyListIsProvided() {
	fileNameList := []string{}

	a, _ := accessWithTestLogging()
	selected := a.filesContainingRSAPrivateKeys(fileNameList)

	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPrivateKeys_ReturnsAnEmptyListIfAListWithANonExistingFileIsProvided() {
	fileNameList := []string{"File that doesn't exist"}

	a, _ := accessWithTestLogging()
	selected := a.filesContainingRSAPrivateKeys(fileNameList)

	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPrivateKeys_ReturnsAnEmptyListIfAListWithAnEmptyFileIsProvided() {
	// Given
	fileName := "Empty file"
	s.createEmptyFile(s.tdir, fileName)
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	a, _ := accessWithTestLogging()
	selected := a.filesContainingRSAPrivateKeys(fileNameList)

	// Then
	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPrivateKeys_ReturnsAnEmptyListIfAListWithAFileThatDoesntContainAnRSAPublicKeyIsProvided() {
	// Given
	fileName := "Empty file"
	s.createFileWithContent(s.tdir, fileName, "not a RSA public key")
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	a, _ := accessWithTestLogging()
	selected := a.filesContainingRSAPrivateKeys(fileNameList)

	// Then
	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPrivateKeys_ReturnsAListWithOneFileNameIfAListWithAFileThatContainsAnRSAPublicKeyIsProvided() {
	// Given
	fileName := "File-with-content"
	s.createFileWithContent(s.tdir, fileName, correctRSASSHPrivateKey)
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	a, _ := accessWithTestLogging()
	selected := a.filesContainingRSAPrivateKeys(fileNameList)

	// Then
	s.Equal(selected, []string{filepath.Join(s.tdir, fileName)})
}

func (s *sshSuite) Test_selectFilesContainingRSAPrivateKeys_ReturnsAListWithSeveralFileNamesThatContainsRSAKey() {
	// Given
	s.createFileWithContent(s.tdir, "key_file1", correctECDSASSHPrivateKey)
	s.createFileWithContent(s.tdir, "key_file2", correctRSASSHPrivateKey)
	s.createEmptyFile(s.tdir, "key_file3")
	s.createFileWithContent(s.tdir, "key_file4", correctRSASSHPrivateKeyOther)

	fileList := s.withDirectory("key_file1", "key_file2", "key_file3", "key_file4", "key_file5")

	// When
	a, _ := accessWithTestLogging()
	selected := a.filesContainingRSAPrivateKeys(fileList)

	// Then
	expected := s.withDirectory("key_file2", "key_file4")
	s.Equal(expected, selected)
}

func (s *sshSuite) Test_checkIfFileContainsAPrivateRSAKey_returnsAnErrorWhenFileDoesntExist() {
	fileName := "a-file-that-doesnt-exist"

	a, _ := accessWithTestLogging()
	_, err := a.checkIfFileContainsAPrivateRSAKey(fileName)

	s.True(os.IsNotExist(err), "Function should generate an error indicating the file doesn`t exist")
}

func (s *sshSuite) Test_checkIfFileContainsAPrivateRSAKey_doesNotReturnErrorWhenFileExists() {
	fileName := "a file that should exist"
	s.createEmptyFile(s.tdir, fileName)

	a, _ := accessWithTestLogging()
	_, err := a.checkIfFileContainsAPrivateRSAKey(filepath.Join(s.tdir, fileName))

	s.Nil(err)
}

func (s *sshSuite) Test_checkIfFileContainsAPrivateRSAKey_ReturnsFalseWhenTheFileContentIsNotAnRSAPrivateKey() {
	fileName := "a file without an RSA Private Key"
	s.createFileWithContent(s.tdir, fileName, "not an RSA Private Key")

	a, _ := accessWithTestLogging()
	b, err := a.checkIfFileContainsAPrivateRSAKey(filepath.Join(s.tdir, fileName))

	s.Nil(err)
	s.False(b)
}

func (s *sshSuite) Test_checkIfFileContainsAPrivateRSAKey_ReturnsTrueWhenTheFileContentIsAnRSAPrivateKey() {
	fileName := "a file with a valid RSA Private Key"
	s.createFileWithContent(s.tdir, fileName, correctRSASSHPrivateKey)

	a, _ := accessWithTestLogging()
	b, err := a.checkIfFileContainsAPrivateRSAKey(filepath.Join(s.tdir, fileName))

	s.Nil(err)
	s.True(b)
}

func (s *sshSuite) Test_removeFileNames_AnUnchangedListIsReturnedWhenNoFileNameToBeRemovedAreGiven() {
	fileNameToDelete := ""

	originalFileNamesList := []string{}
	l := withoutFileName(originalFileNamesList, fileNameToDelete)
	s.Equal(originalFileNamesList, l)

	originalFileNamesList = []string{"file that will not be deleted"}
	l = withoutFileName(originalFileNamesList, fileNameToDelete)
	s.Equal(originalFileNamesList, l)

	originalFileNamesList = []string{"multiple", "files", "that will not", "be deleted"}
	l = withoutFileName(originalFileNamesList, fileNameToDelete)
	s.Equal(originalFileNamesList, l)
}

func (s *sshSuite) Test_removeFileNames_AnUnchangedListIsReturnedWhenThereAreNoCoincidencesWithTheFileNameToBeRemoved() {
	fileNameToDelete := "file to be removed"

	originalFileNamesList := []string{}
	l := withoutFileName(originalFileNamesList, fileNameToDelete)
	s.Equal(originalFileNamesList, l)

	originalFileNamesList = []string{"file that will not be deleted"}
	l = withoutFileName(originalFileNamesList, fileNameToDelete)
	s.Equal(originalFileNamesList, l)

	originalFileNamesList = []string{"multiple", "files", "that will not", "be deleted"}
	l = withoutFileName(originalFileNamesList, fileNameToDelete)
	s.Equal(originalFileNamesList, l)
}

func (s *sshSuite) Test_removeFileName_RemovesTheProvidedFileNameFromAListIfPresent() {
	fileNameToDelete := "coinciding file"

	originalFileNamesList := []string{}
	l := withoutFileName(originalFileNamesList, fileNameToDelete)
	s.Equal([]string{}, l)

	originalFileNamesList = []string{"coinciding file"}
	l = withoutFileName(originalFileNamesList, fileNameToDelete)
	s.Equal([]string{}, l)

	originalFileNamesList = []string{"coinciding file", "not coinciding file"}
	l = withoutFileName(originalFileNamesList, fileNameToDelete)
	s.Equal([]string{"not coinciding file"}, l)
}

func accessWithTestLogging() (*access, *test.Hook) {
	logger, hook := test.NewNullLogger()
	return &access{
		log: logger,
	}, hook
}

func (s *sshSuite) Test_listFilesInHomeSSHDirectory_ReturnsAnEmptyListIfTheDotSSHDirectoryDoesNotExistInTheUsersHomeDirectory() {
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	a, _ := accessWithTestLogging()
	files := a.listFilesInHomeSSHDirectory()

	s.Empty(files)
}

func (s *sshSuite) Test_listFilesInHomeSSHDirectory_ReturnsAnEmptyListIfTheDotSSHDirectoryExistsInTheUsersHomeDirectoryButIsEmpty() {
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	s.Nil(os.Mkdir(path.Join(s.tdir, ".ssh"), 0755))

	a, _ := accessWithTestLogging()
	files := a.listFilesInHomeSSHDirectory()

	s.Empty(files)
}

func (s *sshSuite) Test_listFilesInHomeSSHDirectory_ReturnsAListOfFilesIfTheDotSSHDirectoryExistsInTheUsersHomeDirectoryAndContainsFiles() {
	sshDirectory := path.Join(s.tdir, ".ssh")
	s.Nil(os.Mkdir(sshDirectory, 0755), "it was not possible to create .ssh directory in temporary test directory")
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()

	r := rand.Int()
	expected := []string{"id_rsa.pub", fmt.Sprintf("id_rsa%d", r)}
	for _, f := range expected {
		s.createFileWithContent(sshDirectory, f, "some content")
	}

	expected = transform(expected, func(file string) string {
		return path.Join(sshDirectory, file)
	})

	a, _ := accessWithTestLogging()
	files := a.listFilesInHomeSSHDirectory()

	s.Equal(expected, files)
}

func (s *sshSuite) Test_createPublicKeyEntriesFrom_ReturnsAListOfKeyEntriesFromAllTheProvidedPaths() {
	paths := []string{}
	l := createPublicKeyRepresentationsFrom(paths)
	s.Empty(l)

	paths = []string{"a path"}
	l = createPublicKeyRepresentationsFrom(paths)
	s.Equal([]*publicKeyRepresentation{{"a path"}}, l)

	paths = []string{"a path", "another path"}
	l = createPublicKeyRepresentationsFrom(paths)
	s.Equal([]*publicKeyRepresentation{{"a path"}, {"another path"}}, l)
}

func (s *sshSuite) Test_createPrivateKeyEntriesFrom_ReturnsAListOfKeyEntriesFromAllTheProvidedPaths() {
	paths := []string{}
	l := createPrivateKeyRepresentationsFrom(paths)
	s.Empty(l)

	paths = []string{"a path"}
	l = createPrivateKeyRepresentationsFrom(paths)
	s.Equal([]*privateKeyRepresentation{{"a path"}}, l)

	paths = []string{"a path", "another path"}
	l = createPrivateKeyRepresentationsFrom(paths)
	s.Equal([]*privateKeyRepresentation{{"a path"}, {"another path"}}, l)
}

func (s *sshSuite) Test_privateKeyEntriesFrom_ReturnsAListOfPrivateKeyEntriesFromAllTheProvidedPaths() {
	a, _ := accessWithTestLogging()
	paths := []string{}
	l := a.privateKeyRepresentationsFrom(paths)
	s.Empty(l)

	emptyFile := "Empty-file"
	s.createEmptyFile(s.tdir, emptyFile)
	paths = []string{filepath.Join(s.tdir, emptyFile)}
	l = a.privateKeyRepresentationsFrom(paths)
	s.Empty(l)

	notAnRSAPrivateKeyFile := "Not-an-RSA-private-key-file"
	s.createFileWithContent(s.tdir, notAnRSAPrivateKeyFile, "not a RSA private key")
	paths = []string{
		filepath.Join(s.tdir, emptyFile),
		filepath.Join(s.tdir, notAnRSAPrivateKeyFile),
	}
	l = a.privateKeyRepresentationsFrom(paths)
	s.Empty(l)

	privateRSAKeyFile1 := "File-with-a-private-RSA-key"
	s.createFileWithContent(s.tdir, privateRSAKeyFile1, correctRSASSHPrivateKey)
	privateRSAKeyFile2 := "Another-file-with-another-private-RSA-key"
	s.createFileWithContent(s.tdir, privateRSAKeyFile2, correctRSASSHPrivateKeyOther)

	paths = []string{
		filepath.Join(s.tdir, emptyFile),
		filepath.Join(s.tdir, notAnRSAPrivateKeyFile),
		filepath.Join(s.tdir, privateRSAKeyFile1),
		filepath.Join(s.tdir, privateRSAKeyFile2),
	}

	l = a.privateKeyRepresentationsFrom(paths)
	s.Equal([]*privateKeyRepresentation{
		createPrivateKeyRepresentation(filepath.Join(s.tdir, privateRSAKeyFile1)),
		createPrivateKeyRepresentation(filepath.Join(s.tdir, privateRSAKeyFile2)),
	}, l)
}

func (s *sshSuite) Test_publicKeyEntriesFrom_ReturnsAListOfPublicKeyEntriesFromAllTheProvidedPaths() {
	paths := []string{}
	l := publicKeyRepresentationsFrom(paths)
	s.Empty(l)

	emptyFile := "Empty-file"
	s.createEmptyFile(s.tdir, emptyFile)
	paths = []string{filepath.Join(s.tdir, emptyFile)}
	l = publicKeyRepresentationsFrom(paths)
	s.Empty(l)

	notAnRSAPublicKeyFile := "Not-an-RSA-public-key-file"
	s.createFileWithContent(s.tdir, notAnRSAPublicKeyFile, "not a RSA public key")
	paths = []string{
		filepath.Join(s.tdir, emptyFile),
		filepath.Join(s.tdir, notAnRSAPublicKeyFile),
	}
	l = publicKeyRepresentationsFrom(paths)
	s.Empty(l)

	publicRSAKeyFile1 := "File-with-a-public-RSA-key"
	s.createFileWithContent(s.tdir, publicRSAKeyFile1, "ssh-rsa AAAAA batman@debian")
	publicRSAKeyFile2 := "Another-file-with-another-public-RSA-key"
	s.createFileWithContent(s.tdir, publicRSAKeyFile2, "ssh-rsa NBBBB robin@debian")

	paths = []string{
		filepath.Join(s.tdir, emptyFile),
		filepath.Join(s.tdir, notAnRSAPublicKeyFile),
		filepath.Join(s.tdir, publicRSAKeyFile1),
		filepath.Join(s.tdir, publicRSAKeyFile2),
	}

	l = publicKeyRepresentationsFrom(paths)
	s.Equal([]*publicKeyRepresentation{
		createPublicKeyRepresentation(filepath.Join(s.tdir, publicRSAKeyFile1)),
		createPublicKeyRepresentation(filepath.Join(s.tdir, publicRSAKeyFile2)),
	}, l)
}

func (s *sshSuite) Test_partitionKeyEntries_ReturnsAListOfKeyEntriesWithPublicPrivateAndKeyPairsFromPublicAndPrivateKeyRepresentations() {
	// Both privates and publics are empty
	privates := []*privateKeyRepresentation{}
	publics := []*publicKeyRepresentation{}
	l := partitionKeyEntries(privates, publics)
	s.Empty(l)

	// Only privates and no publics
	privates = []*privateKeyRepresentation{
		createPrivateKeyRepresentation("exclusively"),
		createPrivateKeyRepresentation("privates"),
	}
	publics = []*publicKeyRepresentation{}
	l = partitionKeyEntries(privates, publics)
	s.ElementsMatch([]api.KeyEntry{
		createPrivateKeyRepresentation("exclusively"),
		createPrivateKeyRepresentation("privates"),
	}, l)

	// Only publics and no privates
	privates = []*privateKeyRepresentation{}
	publics = []*publicKeyRepresentation{
		createPublicKeyRepresentation("exclusively.pub"),
		createPublicKeyRepresentation("publics.pub"),
	}
	l = partitionKeyEntries(privates, publics)
	s.ElementsMatch([]api.KeyEntry{
		createPublicKeyRepresentation("exclusively.pub"),
		createPublicKeyRepresentation("publics.pub"),
	}, l)

	// One pair, one lonely public and one lonely private
	privates = []*privateKeyRepresentation{
		createPrivateKeyRepresentation("matching pair"),
		createPrivateKeyRepresentation("lonely private"),
	}
	publics = []*publicKeyRepresentation{
		createPublicKeyRepresentation("matching pair.pub"),
		createPublicKeyRepresentation("lonely public.pub"),
	}
	l = partitionKeyEntries(privates, publics)
	s.ElementsMatch([]api.KeyEntry{
		createKeypairRepresentation(createPrivateKeyRepresentation("matching pair"), createPublicKeyRepresentation("matching pair.pub")),
		createPrivateKeyRepresentation("lonely private"),
		createPublicKeyRepresentation("lonely public.pub"),
	}, l)
}
