package ssh

import (
	"os"
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
