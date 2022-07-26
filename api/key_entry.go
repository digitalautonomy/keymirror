package api

type KeyEntry interface {
	Locations() []string
	PublicKeyLocations() []string
	PrivateKeyLocations() []string
	KeyType() KeyType
}

type PublicKeyEntry interface {
	KeyEntry
	WithDigestContent(func([]byte) []byte) []byte
}
