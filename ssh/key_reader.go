package ssh

import (
	"io/fs"
	"io/ioutil"
)

func listFilesIn(dir string) []string {
	files, _ := ioutil.ReadDir(dir)

	return transform(files, func(f fs.FileInfo) string { return f.Name() })
}

type Key struct {
	t string
}

func parsePublicKey(k string) (Key, bool) {
	if k != "" {
		return Key{"ssh-rsa"}, true
	}

	return Key{}, false
}

// TODO
// - with a list of file names
//   - determine which are RSA public keys
//     - determine if a string has the format of an RSA public key
//        - try to parse string into SSH public key representation
//		  - check if the type identifier is ssh-rsa
//	   - read string from a file, and check if that string is RSA
//   - determine which are RSA private keys
//   - determine if there are duplicates
//   - pair public and private keys to each other based on file name
//   - turn each entry into its internal representation
// - create an internal representation for public, private and key pairs
