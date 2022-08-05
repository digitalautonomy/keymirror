package ssh

import (
	"encoding/base64"
	"regexp"
	"strings"
)

var whitespace = regexp.MustCompile("[[:space:]]+")

type publicKeyParser struct {
	fields []string
}

func newPublicKeyParser(k string) *publicKeyParser {
	p := &publicKeyParser{}
	p.initializeFields(k)
	return p
}

func (p *publicKeyParser) initializeFields(k string) {
	p.fields = whitespace.Split(strings.TrimSpace(k), 3)
}

var commonRSAKeySizes = map[int]int{
	// 1024 bits
	127: 128,
	128: 128,
	129: 128,

	// 2048 bits
	255: 256,
	256: 256,
	257: 256,

	// 3072 bits
	383: 384,
	384: 384,
	385: 384,

	// 4096 bits
	511: 512,
	512: 512,
	513: 512,
}

func canonicalizeRSAKeyLength(p int) int {
	res, ok := commonRSAKeySizes[p]
	if !ok {
		return p
	}
	return res
}

func extractByteLengthFromRSAPublicKey(key []byte) (int, bool) {
	algo, rest, ok := readLengthBytes(key)
	if !ok || string(algo) != rsaAlgorithm {
		return 0, false
	}

	_, rest, ok = readLengthBytes(rest)
	if !ok {
		return 0, false
	}

	product, _, ok := readLengthBytes(rest)
	if !ok {
		return 0, false
	}
	return len(product), true
}

func (p *publicKeyParser) parse() (publicKey, bool) {
	if p.notEnoughFields() {
		return publicKey{}, false
	}

	key, ok := p.parseKey()
	if !ok {
		return publicKey{}, false
	}

	size := 0

	if l, ok2 := extractByteLengthFromRSAPublicKey(key); ok2 {
		size = canonicalizeRSAKeyLength(l) * 8
	} else {
		size = 0
	}

	return publicKey{
		algorithm: p.algorithm(),
		key:       key,
		comment:   p.potentialComment(),
		size:      size,
	}, true
}

func (p *publicKeyParser) notEnoughFields() bool {
	return len(p.fields) == 1
}

func (p *publicKeyParser) algorithm() string {
	return p.fields[0]
}

func (p *publicKeyParser) parseKey() ([]byte, bool) {
	k, e := base64.StdEncoding.DecodeString(p.fields[1])
	return k, e == nil
}

func (p *publicKeyParser) potentialComment() string {
	if p.hasComment() {
		return p.fields[2]
	}
	return ""
}

func (p *publicKeyParser) hasComment() bool {
	return len(p.fields) == 3
}

func parsePublicKey(k string) (publicKey, bool) {
	return newPublicKeyParser(k).parse()
}
