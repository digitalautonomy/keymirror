package ssh

import (
	"io/fs"
	"io/ioutil"
	"regexp"
	"strings"
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

var whitespace = regexp.MustCompile("[[:space:]]+")

func parsePublicKey(k string) (publicKey, bool) {
	fields := whitespace.Split(strings.TrimSpace(k), 3)

	if len(fields) == 1 {
		return publicKey{}, false
	}

	if hasComment(fields) {
		return publicKey{algorithm: fields[0], key: fields[1], comment: fields[2]}, true
	}
	return publicKey{algorithm: fields[0], key: fields[1]}, true
}

func hasComment(fields []string) bool {
	return len(fields) == 3
}

func isRSAPublicKey(k string) bool {
	pub, _ := parsePublicKey(k)
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
