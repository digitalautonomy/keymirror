package api

type KeyAccess interface {
	AllKeys() []KeyEntry
}
