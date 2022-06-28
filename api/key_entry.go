package api

type KeyEntry interface {
	Locations() []string
	PublicKeyLocations() []string
}
