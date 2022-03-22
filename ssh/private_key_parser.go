package ssh

import (
	"bytes"
	"encoding/binary"
	"encoding/pem"
)

type privateKey struct {
	algorithm string
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

func extractPrivateKeyAlgorithm(input []byte) (algorithm string, ok bool) {
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

func createPrivateKeyFrom(input []byte) (privateKey, bool) {
	_, rest, ok1 := readLengthBytes(input)           // reads the ciphername
	_, rest, ok2 := readLengthBytes(rest)            // reads the kdfname
	_, rest, ok3 := readLengthBytes(rest)            // reads the kdf
	numberOfKeys, rest, ok4 := read32BitNumber(rest) // reads the number of keys
	ok5 := numberOfKeys == 1
	_, rest, ok6 := readLengthBytes(rest)  // reads the public key
	value, _, ok7 := readLengthBytes(rest) // reads the private key block
	_, rest, ok8 := extractDummyCheckSum(value)
	algorithm, ok9 := extractPrivateKeyAlgorithm(rest)

	return privateKey{algorithm}, allOK(ok1, ok2, ok3, ok4, ok5, ok6, ok7, ok8, ok9)
}

var privateKeyAuthMagicWithTerminator = []byte("openssh-key-v1\x00")

func parsePrivateKey(pk string) (privateKey, bool) {
	b, _ := pem.Decode([]byte(pk))
	if b == nil || b.Type != "OPENSSH PRIVATE KEY" {
		return privateKey{}, false
	}

	rest, ok := readOpensshPrivateKeyAuthMagic(b.Bytes)
	if !ok {
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
