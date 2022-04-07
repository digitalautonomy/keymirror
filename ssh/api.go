package ssh

import "github.com/digitalautonomy/keymirror/api"

var Access api.KeyAccess = &access{}

type access struct{}

func (*access) AllKeys() []api.KeyEntry {
	files := listFilesInHomeSSHDirectory()
	return partitionKeyEntries(privateKeyRepresentationsFrom(files),
		publicKeyRepresentationsFrom(files))
}
