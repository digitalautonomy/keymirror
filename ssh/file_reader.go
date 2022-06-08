package ssh

import (
	"os"
	"path"
)

func fileContentMatches(fileName string, f predicate[string]) (bool, error) {
	content, e := os.ReadFile(fileName)
	if e != nil {
		return false, e
	}

	return f(string(content)), nil
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
	msg := "found these files in the directory"
	if len(result) == 0 {
		msg = "no files were found in the directory"
	}
	a.log.WithField("ssh files", result).Debug(msg)

	return result
}

func createPublicKeyRepresentationsFrom(input []string) []*publicKeyRepresentation {
	return transform(input, createPublicKeyRepresentation)
}

func createPrivateKeyRepresentationsFrom(input []string) []*privateKeyRepresentation {
	return transform(input, createPrivateKeyRepresentation)
}

func (a *access) privateKeyRepresentationsFrom(input []string) []*privateKeyRepresentation {
	privateKeyFiles := concat(
		a.filesContainingRSAPrivateKeys(input),
		a.filesContainingEd25519PrivateKeys(input),
	)

	return createPrivateKeyRepresentationsFrom(privateKeyFiles)
}

func publicKeyRepresentationsFrom(input []string) []*publicKeyRepresentation {
	publicKeyFiles := concat(
		filesContainingRSAPublicKeys(input),
		filesContainingEd25519PublicKeys(input),
	)
	return createPublicKeyRepresentationsFrom(publicKeyFiles)
}
