package ssh

import (
	"fmt"
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

func (s *sshSuite) Test_selectFilesContainingRSAPublicKeys_ReturnsAListWithSeveralFileNamesThatContainsRSAKey() {
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

func (s *sshSuite) Test_selectFilesContainingRSAPrivateKeys_ReturnsAnEmptyListIfAnEmptyListIsProvided() {
	fileNameList := []string{}

	selected := selectFilesContainingRSAPrivateKeys(fileNameList)

	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPrivateKeys_ReturnsAnEmptyListIfAListWithANonExistingFileIsProvided() {
	fileNameList := []string{"File that doesn't exist"}

	selected := selectFilesContainingRSAPrivateKeys(fileNameList)

	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPrivateKeys_ReturnsAnEmptyListIfAListWithAnEmptyFileIsProvided() {
	// Given
	fileName := "Empty file"
	s.createEmptyFile(s.tdir, fileName)
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	selected := selectFilesContainingRSAPrivateKeys(fileNameList)

	// Then
	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPrivateKeys_ReturnsAnEmptyListIfAListWithAFileThatDoesntContainAnRSAPublicKeyIsProvided() {
	// Given
	fileName := "Empty file"
	s.createFileWithContent(s.tdir, fileName, "not a RSA public key")
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	selected := selectFilesContainingRSAPrivateKeys(fileNameList)

	// Then
	s.Empty(selected)
}

func (s *sshSuite) Test_selectFilesContainingRSAPrivateKeys_ReturnsAListWithOneFileNameIfAListWithAFileThatContainsAnRSAPublicKeyIsProvided() {
	// Given
	fileName := "File-with-content"
	s.createFileWithContent(s.tdir, fileName, correctRSASSHPrivateKey)
	fileNameList := []string{filepath.Join(s.tdir, fileName)}

	// When
	selected := selectFilesContainingRSAPrivateKeys(fileNameList)

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
	selected := selectFilesContainingRSAPrivateKeys(fileList)

	// Then
	expected := s.withDirectory("key_file2", "key_file4")
	s.Equal(expected, selected)
}

func (s *sshSuite) Test_checkIfFileContainsAPrivateRSAKey_returnsAnErrorWhenFileDoesntExist() {
	fileName := "a-file-that-doesnt-exist"

	_, err := checkIfFileContainsAPrivateRSAKey(fileName)

	s.True(os.IsNotExist(err), "Function should generate an error indicating the file doesn`t exist")
}

func (s *sshSuite) Test_checkIfFileContainsAPrivateRSAKey_doesNotReturnErrorWhenFileExists() {
	fileName := "a file that should exist"
	s.createEmptyFile(s.tdir, fileName)

	_, err := checkIfFileContainsAPrivateRSAKey(filepath.Join(s.tdir, fileName))

	s.Nil(err)
}

func (s *sshSuite) Test_checkIfFileContainsAPrivateRSAKey_ReturnsFalseWhenTheFileContentIsNotAnRSAPrivateKey() {
	fileName := "a file without an RSA Private Key"
	s.createFileWithContent(s.tdir, fileName, "not an RSA Private Key")

	b, err := checkIfFileContainsAPrivateRSAKey(filepath.Join(s.tdir, fileName))

	s.Nil(err)
	s.False(b)
}

func (s *sshSuite) Test_checkIfFileContainsAPrivateRSAKey_ReturnsTrueWhenTheFileContentIsAnRSAPrivateKey() {
	fileName := "a file with a valid RSA Private Key"
	s.createFileWithContent(s.tdir, fileName, correctRSASSHPrivateKey)

	b, err := checkIfFileContainsAPrivateRSAKey(filepath.Join(s.tdir, fileName))

	s.Nil(err)
	s.True(b)
}

func (s *sshSuite) Test_removePubSuffixFromFileName_AnEmptyListReturnsAnEmptyList() {
	files := []string{}

	l := removePubSuffixFromFileNamesList(files)

	s.Empty(l)
}

func (s *sshSuite) Test_removePubSuffixFromFileName_AListWithOneFileNameWithoutPubSuffixIsNotModified() {
	files := []string{"key1.bla"}

	l := removePubSuffixFromFileNamesList(files)

	s.Equal([]string{"key1.bla"}, l)
}

func (s *sshSuite) Test_removePubSuffixFromFileName_RemovesPubSuffixFromFileNameFromAListWithOneFileNameWithPubSuffix() {
	files := []string{"key1.pub"}

	l := removePubSuffixFromFileNamesList(files)

	s.Equal([]string{"key1"}, l)
}

func (s *sshSuite) Test_removePubSuffixFromFileName_RemovesPubSuffixFromFileNamesIfTheyHavePubSuffix() {
	files := []string{"key1.pub", "key2.bla", "key3", "key4.pub"}

	l := removePubSuffixFromFileNamesList(files)

	s.Equal([]string{"key1", "key2.bla", "key3", "key4"}, l)
}

func (s *sshSuite) Test_findKeyPairsBasedOnFileName_AnEmptyListIsReturnedWhenTwoEmptyListsAreGiven() {
	privateFiles := []string{}
	publicFiles := []string{}

	l := findKeyPairsBasedOnFileName(privateFiles, publicFiles)

	s.Equal([]string{}, l)
}

func (s *sshSuite) Test_findKeyPairsBasedOnFileName_AnEmptyListIsReturnedWhenOneOfTheGivenListsIsEmpty() {
	privateFiles := []string{"key1"}
	publicFiles := []string{}
	l := findKeyPairsBasedOnFileName(privateFiles, publicFiles)
	s.Equal([]string{}, l)

	privateFiles = []string{}
	publicFiles = []string{"key2"}
	l = findKeyPairsBasedOnFileName(privateFiles, publicFiles)
	s.Equal([]string{}, l)
}

func (s *sshSuite) Test_findKeyPairsBasedOnFileName_AnEmptyListIsReturnedWhenTwoListsWithoutFileNamesCoincidencesAreGiven() {
	privateFiles := []string{"key1"}
	publicFiles := []string{"key2"}

	l := findKeyPairsBasedOnFileName(privateFiles, publicFiles)

	s.Equal([]string{}, l)
}

func (s *sshSuite) Test_findKeyPairsBasedOnFileName_AListContainingTheOnlyCoincidingFileNameIsReturnedWhenTwoListsWithTheSameFileNameAreGiven() {
	privateFiles := []string{"key3"}
	publicFiles := []string{"key3"}

	l := findKeyPairsBasedOnFileName(privateFiles, publicFiles)

	s.Equal([]string{"key3"}, l)
}

func (s *sshSuite) Test_findKeyPairsBasedOnFileName_AListContainingMultipleCoincidingFileNamesIsReturnedWhenTwoListsWithMultipleCoincidingFileNamesAreGiven() {
	privateFiles := []string{"key1", "key2", "key3", "key4"}
	publicFiles := []string{"key2", "key3"}

	l := findKeyPairsBasedOnFileName(privateFiles, publicFiles)

	s.Equal([]string{"key2", "key3"}, l)
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

func (s *sshSuite) Test_removeFileNames_RemovesTheProvidedFileNamesFromAListIfPresent() {
	fileNamesToDelete := []string{"coinciding", "files"}

	originalFileNamesList := []string{}
	l := removeFileNames(originalFileNamesList, fileNamesToDelete)
	s.Equal([]string{}, l)

	originalFileNamesList = []string{"coinciding", "files"}
	l = removeFileNames(originalFileNamesList, fileNamesToDelete)
	s.Equal([]string{}, l)

	originalFileNamesList = []string{"coinciding", "files", "not coinciding file"}
	l = removeFileNames(originalFileNamesList, fileNamesToDelete)
	s.Equal([]string{"not coinciding file"}, l)
}

func (s *sshSuite) Test_listFilesInHomeSSHDirectory_ReturnsAnEmptyListIfTheDotSSHDirectoryDoesNotExistInTheUsersHomeDirectory() {
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	files := listFilesInHomeSSHDirectory()

	s.Empty(files)
}

func (s *sshSuite) Test_listFilesInHomeSSHDirectory_ReturnsAnEmptyListIfTheDotSSHDirectoryExistsInTheUsersHomeDirectoryButIsEmpty() {
	defer gostub.New().SetEnv("HOME", s.tdir).Reset()
	s.Nil(os.Mkdir(path.Join(s.tdir, ".ssh"), 0755))

	files := listFilesInHomeSSHDirectory()

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

	files := listFilesInHomeSSHDirectory()

	s.Equal(expected, files)
}

func (s *sshSuite) Test_createPublicKeyEntriesFrom_ReturnsAListOfKeyEntriesFromAllTheProvidedPaths() {
	paths := []string{}
	l := createPublicKeyEntriesFrom(paths)
	s.Empty(l)

	paths = []string{"a path"}
	l = createPublicKeyEntriesFrom(paths)
	s.Equal([]KeyEntry{&publicKeyRepresentation{"a path"}}, l)

	paths = []string{"a path", "another path"}
	l = createPublicKeyEntriesFrom(paths)
	s.Equal([]KeyEntry{&publicKeyRepresentation{"a path"}, &publicKeyRepresentation{"another path"}}, l)
}

func (s *sshSuite) Test_createPrivateKeyEntriesFrom_ReturnsAListOfKeyEntriesFromAllTheProvidedPaths() {
	paths := []string{}
	l := createPrivateKeyEntriesFrom(paths)
	s.Empty(l)

	paths = []string{"a path"}
	l = createPrivateKeyEntriesFrom(paths)
	s.Equal([]KeyEntry{&privateKeyRepresentation{"a path"}}, l)

	paths = []string{"a path", "another path"}
	l = createPrivateKeyEntriesFrom(paths)
	s.Equal([]KeyEntry{&privateKeyRepresentation{"a path"}, &privateKeyRepresentation{"another path"}}, l)
}

func (s *sshSuite) Test_privateKeyEntriesFrom_ReturnsAListOfPrivateKeyEntriesFromAllTheProvidedPaths() {
	paths := []string{}
	l := privateKeyEntriesFrom(paths)
	s.Empty(l)

	emptyFile := "Empty-file"
	s.createEmptyFile(s.tdir, emptyFile)
	paths = []string{filepath.Join(s.tdir, emptyFile)}
	l = privateKeyEntriesFrom(paths)
	s.Empty(l)

	notAnRSAPrivateKeyFile := "Not-an-RSA-private-key-file"
	s.createFileWithContent(s.tdir, notAnRSAPrivateKeyFile, "not a RSA private key")
	paths = []string{
		filepath.Join(s.tdir, emptyFile),
		filepath.Join(s.tdir, notAnRSAPrivateKeyFile),
	}
	l = privateKeyEntriesFrom(paths)
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

	l = privateKeyEntriesFrom(paths)
	s.Equal([]KeyEntry{
		createPrivateKeyRepresentation(filepath.Join(s.tdir, privateRSAKeyFile1)),
		createPrivateKeyRepresentation(filepath.Join(s.tdir, privateRSAKeyFile2)),
	}, l)
}

func (s *sshSuite) Test_publicKeyEntriesFrom_ReturnsAListOfPublicKeyEntriesFromAllTheProvidedPaths() {
	paths := []string{}
	l := publicKeyEntriesFrom(paths)
	s.Empty(l)

	emptyFile := "Empty-file"
	s.createEmptyFile(s.tdir, emptyFile)
	paths = []string{filepath.Join(s.tdir, emptyFile)}
	l = publicKeyEntriesFrom(paths)
	s.Empty(l)

	notAnRSAPublicKeyFile := "Not-an-RSA-public-key-file"
	s.createFileWithContent(s.tdir, notAnRSAPublicKeyFile, "not a RSA public key")
	paths = []string{
		filepath.Join(s.tdir, emptyFile),
		filepath.Join(s.tdir, notAnRSAPublicKeyFile),
	}
	l = publicKeyEntriesFrom(paths)
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

	l = publicKeyEntriesFrom(paths)
	s.Equal([]KeyEntry{
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
	s.Equal([]KeyEntry{
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
	s.Equal([]KeyEntry{
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
	s.Equal([]KeyEntry{
		createKeypairRepresentation(createPrivateKeyRepresentation("matching pair"), createPublicKeyRepresentation("matching pair.pub")),
		createPrivateKeyRepresentation("lonely private"),
		createPublicKeyRepresentation("lonely public.pub"),
	}, l)
}
