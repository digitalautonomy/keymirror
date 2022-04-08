package ssh

import (
	"os"
	"path"
)

func checkIfFileContainsAPublicRSAKey(fileName string) (bool, error) {
	return checkIfFileContainsASpecificValue(fileName, isRSAPublicKey)
}

func (a *access) checkIfFileContainsAPrivateRSAKey(fileName string) (bool, error) {
	return checkIfFileContainsASpecificValue(fileName, a.isRSAPrivateKey)
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

func (a *access) filesContainingRSAPrivateKeys(fileNameList []string) []string {
	a.log.WithField("file names to check", fileNameList).Trace("filesContainingRSAPrivateKeys()")
	result := filter(fileNameList, loggingErrors(a.log, "an error happened while checking if a file contains a private key", a.checkIfFileContainsAPrivateRSAKey))
	a.log.WithField("private key files", result).Debug("we found these RSA private key files")
	return result
}

func withoutFileName(targetFileNamesList []string, fileNameToDelete string) []string {
	return filter(targetFileNamesList, not(isEqualTo(fileNameToDelete)))
}

func (a *access) listFilesInHomeSSHDirectory() []string {
	sshDirectory := path.Join(os.Getenv("HOME"), ".ssh")
	a.log.WithField("ssh directory", sshDirectory).Debug("listing files in users .ssh home directory")
	result := transform(listFilesIn(sshDirectory), func(file string) string {
		return path.Join(sshDirectory, file)
	})
	a.log.WithField("ssh files", result).Debug("found these files in the directory")
	return result
}

func createPublicKeyRepresentationsFrom(input []string) []*publicKeyRepresentation {
	return transform(input, createPublicKeyRepresentation)
}

func createPrivateKeyRepresentationsFrom(input []string) []*privateKeyRepresentation {
	return transform(input, createPrivateKeyRepresentation)
}

func (a *access) privateKeyRepresentationsFrom(input []string) []*privateKeyRepresentation {
	return createPrivateKeyRepresentationsFrom(a.filesContainingRSAPrivateKeys(input))
}

func publicKeyRepresentationsFrom(input []string) []*publicKeyRepresentation {
	return createPublicKeyRepresentationsFrom(filesContainingRSAPublicKeys(input))
}
