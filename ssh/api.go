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
	// output: the list of Private Keys file names -> selectFilesContainingRSAPublicKeys returns as list of strings

	// input: the list of files names of files in ~/.ssh directory as list of strings
	// output: the list of Public Keys file names -> selectFilesContainingRSAPrivateKeys returns as list of strings

	// input: a list of strings corresponding to the file names retrieved in the 2 steps before
	// output: a list of KeyEntry find pairs, lonely public, lonely private and return them -> defineKeyTypesFrom

	return nil
}

// in th egui

//func bla() {
//	keysToList := ssh.Access.AllKeys()
//	// display keys to listr
//}
