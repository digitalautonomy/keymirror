package ssh

import (
	"io/fs"
	"io/ioutil"
)

func listFilesIn(dir string) []string {
	files, _ := ioutil.ReadDir(dir)

	return transform(files, (fs.FileInfo).Name)
}

type publicKey struct {
	location  string
	algorithm string
	key       string
	comment   string
}

func (k *publicKey) isAlgorithm(algo string) bool {
	return k.algorithm == algo
}
