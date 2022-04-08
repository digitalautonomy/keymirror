package ssh

import (
	"github.com/digitalautonomy/keymirror/api"
	"github.com/sirupsen/logrus"
)

func Access(l logrus.FieldLogger) api.KeyAccess {
	return &access{log: l.WithField("component", "ssh")}
}

type access struct {
	log logrus.Ext1FieldLogger
}

func (a *access) AllKeys() []api.KeyEntry {
	files := a.listFilesInHomeSSHDirectory()
	return partitionKeyEntries(a.privateKeyRepresentationsFrom(files),
		publicKeyRepresentationsFrom(files))
}
