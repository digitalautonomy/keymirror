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

	// find if the .ssh directory exists in home directory (~/.ssh)
	// if it does not exist return an empty list of fileNames
	// if it exists in read the list of files it contains

	// input .ssh directory
	// output: list of fileNames

	// input: the list of file names in ~/.ssh directory as list of strings

	// detectPrivateKeys
	// input: a list of file names
	// output: a list of private key representations

	// detectPublicKeys
	// input: a list of file names
	// output: a list of public key representations

	// input: the list of public representations, and the list of private representations
	// output: one list of public/private/key pair representations

	return nil
}

// in th egui

//func bla() {
//	keysToList := ssh.Access.AllKeys()
//	// display keys to listr
//}
