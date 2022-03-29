package ssh

import (
	"os"
	"path"
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

func filesContainingRSAPublicKeys(fileNameList []string) []string {
	return filter(fileNameList, ignoringErrors(checkIfFileContainsAPublicRSAKey))
}

func filesContainingRSAPrivateKeys(fileNameList []string) []string {
	return filter(fileNameList, ignoringErrors(checkIfFileContainsAPrivateRSAKey))
}

func withoutFileName(targetFileNamesList []string, fileNameToDelete string) []string {
	return filter(targetFileNamesList, not(isEqualTo(fileNameToDelete)))
}

func listFilesInHomeSSHDirectory() []string {
	sshDirectory := path.Join(os.Getenv("HOME"), ".ssh")
	return transform(listFilesIn(sshDirectory), func(file string) string {
		return path.Join(sshDirectory, file)
	})

}

func createPublicKeyRepresentationsFrom(input []string) []*publicKeyRepresentation {
	return transform(input, createPublicKeyRepresentation)
}

func createPrivateKeyRepresentationsFrom(input []string) []*privateKeyRepresentation {
	return transform(input, createPrivateKeyRepresentation)
}

func privateKeyRepresentationsFrom(input []string) []*privateKeyRepresentation {
	return createPrivateKeyRepresentationsFrom(filesContainingRSAPrivateKeys(input))
}

func publicKeyRepresentationsFrom(input []string) []*publicKeyRepresentation {
	return createPublicKeyRepresentationsFrom(filesContainingRSAPublicKeys(input))
}
