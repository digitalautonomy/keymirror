package ssh

import (
	"os"
	"path"
	"strings"
)

func checkIfFileContainsAPublicRSAKey(fileName string) (bool, error) {
	return checkIfFileContainsASpecificValue(fileName, isRSAPublicKey)
}

func checkIfFileContainsAPrivateRSAKey(fileName string) (bool, error) {
	return checkIfFileContainsASpecificValue(fileName, isRSAPrivateKey)
}

func checkIfFileContainsASpecificValue(fileName string, f predicate[string]) (bool, error) {
	content, e := os.ReadFile(fileName)
	if e != nil {
		return false, e
	}

	return f(string(content)), nil
}

func selectFilesContainingRSAPublicKeys(fileNameList []string) []string {
	return filter(fileNameList, ignoringErrors(checkIfFileContainsAPublicRSAKey))
}

func selectFilesContainingRSAPrivateKeys(fileNameList []string) []string {
	return filter(fileNameList, ignoringErrors(checkIfFileContainsAPrivateRSAKey))
}

func removePubSuffixFromFileName(s string) string {
	return strings.TrimSuffix(s, ".pub")
}

func removePubSuffixFromFileNamesList(files []string) []string {
	return transform(files, removePubSuffixFromFileName)
}

func findKeyPairsBasedOnFileName(privateFiles, publicFiles []string) []string {
	return filter(privateFiles, existsIn(publicFiles))
}

func withoutFileName(targetFileNamesList []string, fileNameToDelete string) []string {
	return filter(targetFileNamesList, not(isEqualTo(fileNameToDelete)))
}

func removeFileNames(targetFileNamesList, fileNameToDelete []string) []string {
	return foldLeft(fileNameToDelete, targetFileNamesList, withoutFileName)
}

func listFilesInHomeSSHDirectory() []string {
	sshDirectory := path.Join(os.Getenv("HOME"), ".ssh")
	return listFilesIn(sshDirectory)
}

//func defineKeyTypesFrom(privateKeyFileNames, publicKeyFileNames []string) map[string]string {
//	if isEmptySlice(privateKeyFileNames) {
//		return map[string]string{}
//	}
//	return map[string]string{"privateKeyFile1": "private"}
//}
