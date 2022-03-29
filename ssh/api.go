package ssh

type KeyEntry interface {
	Locations() []string
}

type KeyAccess interface {
	AllKeys() []KeyEntry
}

var Access KeyAccess = &access{}

type access struct{}

func (*access) AllKeys() []KeyEntry {
	files := listFilesInHomeSSHDirectory()
	return partitionKeyEntries(privateKeyRepresentationsFrom(files),
		publicKeyRepresentationsFrom(files))
}
