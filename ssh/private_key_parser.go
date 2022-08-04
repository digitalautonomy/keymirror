package ssh

import (
	"bytes"
	"encoding/binary"
	"encoding/pem"
)

type privateKey struct {
	path              string
	algorithm         string
	passwordProtected bool
	size              int
}

func (k *privateKey) isAlgorithm(algo string) bool {
	return k.algorithm == algo
}

func readBytes(input []byte, n int) (read []byte, rest []byte, ok bool) {
	if len(input) < n {
		return
	}

	return input[:n], input[n:], true
}

func read32BitNumber(input []byte) (value uint32, rest []byte, ok bool) {
	read, rest, ok := readBytes(input, 4)
	if !ok {
		return
	}

	value = binary.BigEndian.Uint32(read)
	return value, rest, ok
}

func readLengthBytes(input []byte) (value []byte, rest []byte, ok bool) {
	l, rest1, ok1 := read32BitNumber(input)
	if !ok1 {
		return nil, nil, false
	}
	return readBytes(rest1, int(l))
}

func extractKeyAlgorithm(input []byte) (algorithm string, ok bool) {
	a, _, k := readLengthBytes(input)
	return string(a), k
}

func extractDummyCheckSum(input []byte) (checksum []byte, rest []byte, ok bool) {
	return readBytes(input, 8)
}

func allOK(vals ...bool) bool {
	if len(vals) == 1 {
		return vals[0]
	}
	return vals[0] && allOK(vals[1:]...)
}

func hasNoAlgorithm(v []byte) bool {
	return string(v) == "none"
}

func createPrivateKeyFrom(input []byte) (privateKey, bool) {
	cipherName, rest, ok1 := readLengthBytes(input)  // reads the ciphername
	_, rest, ok2 := readLengthBytes(rest)            // reads the kdfname
	_, rest, ok3 := readLengthBytes(rest)            // reads the kdf
	numberOfKeys, rest, ok4 := read32BitNumber(rest) // reads the number of keys
	ok5 := numberOfKeys == 1
	pubValue, rest, ok6 := readLengthBytes(rest) // reads the public key
	privValue, _, ok7 := readLengthBytes(rest)   // reads the private key block

	l, _ := extractByteLengthFromRSAPublicKey(pubValue)
	size := canonicalizeRSAKeyLength(l) * 8

	if hasNoAlgorithm(cipherName) {
		_, rest, ok8 := extractDummyCheckSum(privValue)
		algorithm, ok9 := extractKeyAlgorithm(rest)
		return privateKey{
			algorithm:         algorithm,
			passwordProtected: false,
			size:              size,
		}, allOK(ok1, ok2, ok3, ok4, ok5, ok6, ok7, ok8, ok9)
	}

	algorithm, ok8 := extractKeyAlgorithm(pubValue)
	return privateKey{
		algorithm:         algorithm,
		passwordProtected: true,
		size:              size,
	}, allOK(ok1, ok2, ok3, ok4, ok5, ok6, ok7, ok8)
}

var privateKeyAuthMagicWithTerminator = []byte("openssh-key-v1\x00")

func (a *access) parsePrivateKey(pk string) (privateKey, bool) {
	b, _ := pem.Decode([]byte(pk))
	if b == nil || b.Type != "OPENSSH PRIVATE KEY" {
		if b == nil {
			a.log.Error("PEM decoding of RSA private key failed")
		} else {
			a.log.WithField("pem type", b.Type).Error("Incorrect PEM type for RSA private key")
		}
		return privateKey{}, false
	}

	rest, ok := readOpensshPrivateKeyAuthMagic(b.Bytes)
	if !ok {
		a.log.Error("Incorrect Openssh private key magic string")
		return privateKey{}, false
	}

	return createPrivateKeyFrom(rest)
}

func readOpensshPrivateKeyAuthMagic(input []byte) (rest []byte, ok bool) {
	if !bytes.HasPrefix(input, privateKeyAuthMagicWithTerminator) {
		return nil, false
	}

	return bytes.TrimPrefix(input, privateKeyAuthMagicWithTerminator), true
}
