package ssh

func (a *access) checkIfFileContainsAPrivateEd25519Key(fileName string) (bool, error) {
	return fileContentMatches(fileName, a.isEd25519PrivateKey)
}

func (a *access) filesContainingEd25519PrivateKeys(fileNameList []string) []string {
	result := filter(fileNameList, ignoringErrors(a.checkIfFileContainsAPrivateEd25519Key))
	return result
}

func ed25519PublicKeyFrom(fileNameList []string) []*publicKey {
	return filter(transform(fileNameList, publicKeyFromFile), both(not(isNil[publicKey]), (*publicKey).isEd25519))
}

func (a *access) ed25519PrivateKeyFrom(fileNameList []string) []*privateKey {
	return filter(transform(fileNameList, a.privateKeyFromFile), both(not(isNil[privateKey]), (*privateKey).isEd25519))
}
