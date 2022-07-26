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

func (p *publicKeyParser) parse() (publicKey, bool) {
	if p.notEnoughFields() {
		return publicKey{}, false
	}

	key, ok := p.parseKey()
	if !ok {
		return publicKey{}, false
	}

	return publicKey{
		algorithm: p.algorithm(),
		key:       key,
		comment:   p.potentialComment(),
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
