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

func createPublicKeyRepresentationsFromPublicKeys(input []*publicKey) []*publicKeyRepresentation {
	return transform(input, createPublicKeyRepresentationFromPublicKey)
}

func createPrivateKeyRepresentationFromPrivateKeys(input []*privateKey) []*privateKeyRepresentation {
	return transform(input, createPrivateKeyRepresentationFromPrivateKey)
}

func (a *access) privateKeyRepresentationsFrom(input []string) []*privateKeyRepresentation {
	rsaKeys := a.rsaPrivateKeyFrom(input)
	rsaKeyRepresentations := createPrivateKeyRepresentationFromPrivateKeys(rsaKeys)

	ed25519Keys := a.ed25519PrivateKeyFrom(input)
	ed25519KeyRepresentations := createPrivateKeyRepresentationFromPrivateKeys(ed25519Keys)

	return concat(
		rsaKeyRepresentations,
		ed25519KeyRepresentations,
	)
}

func publicKeyRepresentationsFrom(input []string) []*publicKeyRepresentation {
	rsaKeys := rsaPublicKeysFrom(input)
	rsaKeyRepresentations := createPublicKeyRepresentationsFromPublicKeys(rsaKeys)

	ed25519Keys := ed25519PublicKeyFrom(input)
	ed25519KeyRepresentations := createPublicKeyRepresentationsFromPublicKeys(ed25519Keys)

	return concat(
		rsaKeyRepresentations,
		ed25519KeyRepresentations)
}
