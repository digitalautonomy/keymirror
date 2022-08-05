package api

type KeyEntry interface {
	Locations() []string
	PublicKeyLocations() []string
	PrivateKeyLocations() []string
	KeyType() KeyType
	Size() int
	Algorithm() Algorithm
}

type PublicKeyEntry interface {
	KeyEntry
	WithDigestContent(func([]byte) []byte) []byte
	UserID() string
}

type PrivateKeyEntry interface {
	KeyEntry
	IsPasswordProtected() bool
}
