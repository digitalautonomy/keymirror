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

// TODO
// - with a list of file names
//   - determine which are RSA public keys ✓
//     - determine if a string has the format of an RSA public key ✓
//        - try to parse string into SSH public key representation ✓
//		  - check if the type identifier is ssh-rsa ✓
//	   - read string from a file, and check if that string is RSA ✓
//   - determine which are RSA private keys ✓
//     - determine if a string has the format of an RSA private key ✓
//       - find or create a basic PEM reader that returns two things for a file - the "tag" for example OPENSSH PRIVATE KEY or RSA PRIVATE KEY, and the main content
//         either in base64 or in binary, where the base64 has been unpacked. PEM supports headers, and if the library supports that, fine, but it is not necessary now ✓
//       - decode enough of the binary to find ONLY the place that shows that this binary is an SSH private key, and that the type of the key corresponds to RSA ✓
//         - binary has to begin with openssh-key-v1 0x00 ✓
//         - we dont care at this point if its encrypted or not. we need to read length and jump past 4 fields. ✓
//         - we need to read the length for the public key and jump past it. More or less we need to get to the 7th field, which contain the private key keytype ✓
//	   - read string from a file, and check if that string is RSA private key ✓
//   - pair public and private keys to each other based on file name
//   - turn each entry into its internal representation
// - go through all files in .ssh, detect if they are public or private keys, and
//   create a list for that. then, create a new list with unified content
// - create an internal representation for public, private and key pairs (waiting for GUI needs)

// - file names
// - indication if public, private or key pair
