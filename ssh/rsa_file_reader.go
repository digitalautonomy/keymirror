package ssh

import "os"

func checkIfFileContainsAPublicRSAKey(fileName string) (bool, error) {
	return fileContentMatches(fileName, isRSAPublicKey)
}

func (a *access) checkIfFileContainsAPrivateRSAKey(fileName string) (bool, error) {
	return fileContentMatches(fileName, a.isRSAPrivateKey)
}

func filesContainingRSAPublicKeys(fileNameList []string) []string {
	return filter(fileNameList, ignoringErrors(checkIfFileContainsAPublicRSAKey))
}

func (a *access) privateKeyFromFile(fileName string) *privateKey {
	content, e := os.ReadFile(fileName)
	if e != nil {
		return nil
	}

	priv, ok := a.parsePrivateKey(string(content))
	if !ok {
		return nil
	}

	priv.path = fileName
	return &priv
}

func publicKeyFromFile(fileName string) *publicKey {
	content, e := os.ReadFile(fileName)
	if e != nil {
		return nil
	}

	pub, ok := parsePublicKey(string(content))
	if !ok {
		return nil
	}

	pub.location = fileName
	return &pub
}

func (a *access) rsaPrivateKeyFrom(fileNameList []string) []*privateKey {
	return filter(transform(fileNameList, a.privateKeyFromFile), both(not(isNil[privateKey]), (*privateKey).isRSA))
}

func rsaPublicKeysFrom(fileNameList []string) []*publicKey {
	return filter(transform(fileNameList, publicKeyFromFile), both(not(isNil[publicKey]), (*publicKey).isRSA))
}

func (a *access) filesContainingRSAPrivateKeys(fileNameList []string) []string {
	a.log.WithField("file names to check", fileNameList).Trace("filesContainingRSAPrivateKeys()")
	result := filter(fileNameList, loggingErrors(a.log, "an error happened while checking if a file contains a private key", a.checkIfFileContainsAPrivateRSAKey))
	a.log.WithField("private key files", result).Debug("we found these RSA private key files")
	return result
}
