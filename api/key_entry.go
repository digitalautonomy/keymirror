package api

type KeyEntry interface {
	Locations() []string
	PublicKeyLocations() []string
	PrivateKeyLocations() []string
	KeyType() KeyType
	Size() int
}

type PublicKeyEntry interface {
	KeyEntry
	WithDigestContent(func([]byte) []byte) []byte
}

type PrivateKeyEntry interface {
	KeyEntry
	IsPasswordProtected() bool
}
