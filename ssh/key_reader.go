package ssh

import (
	"io/fs"
	"io/ioutil"
)

func listFilesIn(dir string) []string {
	files, _ := ioutil.ReadDir(dir)

	return transform(files, func(f fs.FileInfo) string { return f.Name() })
}

type publicKey struct {
	algorithm string
	key       string
	comment   string
}

const rsaAlgorithm = "ssh-rsa"

func (k *publicKey) isAlgorithm(algo string) bool {
	return k.algorithm == algo
}

func (k *publicKey) isRSA() bool {
	return k.isAlgorithm(rsaAlgorithm)
}

func isRSAPublicKey(k string) bool {
	pub, ok := parsePublicKey(k)
	if !ok {
		return false
	}
	return pub.isRSA()
}

func isRSAPrivateKey(pk string) bool {
	priv, ok := parsePrivateKey(pk)
	if !ok {
		return false
	}
	return priv.isRSA()
}
