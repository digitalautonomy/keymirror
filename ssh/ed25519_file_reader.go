package ssh

func (a *access) checkIfFileContainsAPrivateEd25519Key(fileName string) (bool, error) {
	return fileContentMatches(fileName, a.isEd25519PrivateKey)
}

func (a *access) filesContainingEd25519PrivateKeys(fileNameList []string) []string {
	result := filter(fileNameList, ignoringErrors(a.checkIfFileContainsAPrivateEd25519Key))
	return result
}

func checkIfFileContainsAPublicEd25519Key(fileName string) (bool, error) {
	return fileContentMatches(fileName, isEd25519PublicKey)
}

func filesContainingEd25519PublicKeys(fileNameList []string) []string {
	return filter(fileNameList, ignoringErrors(checkIfFileContainsAPublicEd25519Key))
}
