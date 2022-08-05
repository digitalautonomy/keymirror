package api

type Algorithm interface {
	HasKeySize() bool
	Name() string
}

type algorithm struct {
	hasKeySize bool
	name       string
}

func (a *algorithm) HasKeySize() bool {
	return a.hasKeySize
}

func (a *algorithm) Name() string {
	return a.name
}

var RSA Algorithm = &algorithm{hasKeySize: true, name: "RSA"}
var Ed25519 Algorithm = &algorithm{hasKeySize: false, name: "Ed25519"}
var DSA Algorithm = &algorithm{hasKeySize: false, name: "DSA"}
