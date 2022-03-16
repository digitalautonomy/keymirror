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

// TODO
// - with a list of file names
//   - determine which are RSA public keys
//     - determine if a string has the format of an RSA public key ✓
//        - try to parse string into SSH public key representation ✓
//		  - check if the type identifier is ssh-rsa ✓
//	   - read string from a file, and check if that string is RSA
//   - determine which are RSA private keys
//   - determine if there are duplicates
//   - pair public and private keys to each other based on file name
//   - turn each entry into its internal representation
// - create an internal representation for public, private and key pairs
