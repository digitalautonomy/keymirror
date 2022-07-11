package api

type KeyEntry interface {
	Locations() []string
	PublicKeyLocations() []string
	PrivateKeyLocations() []string
	KeyType() KeyType
}
