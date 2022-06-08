package ssh

func checkIfFileContainsAPublicRSAKey(fileName string) (bool, error) {
	return fileContentMatches(fileName, isRSAPublicKey)
}

func (a *access) checkIfFileContainsAPrivateRSAKey(fileName string) (bool, error) {
	return fileContentMatches(fileName, a.isRSAPrivateKey)
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
