package ssh

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func checkIfFileContainsAPublicRSAKey(fileName string) (bool, error) {
	return checkIfFileContainsASpecificValue(fileName, isRSAPublicKey)
}

func checkIfFileContainsAPrivateRSAKey(fileName string) (bool, error) {
	return checkIfFileContainsASpecificValue(fileName, isRSAPrivateKey)
}

func checkIfFileContainsASpecificValue(fileName string, f predicate[string]) (bool, error) {
	content, e := os.ReadFile(fileName)
	if e != nil {
		return false, e
	}

	return f(string(content)), nil
}

func selectFilesContainingRSAPublicKeys(fileNameList []string) []string {
	return filter(fileNameList, ignoringErrors(checkIfFileContainsAPublicRSAKey))
}

func selectFilesContainingRSAPrivateKeys(fileNameList []string) []string {
	return filter(fileNameList, ignoringErrors(checkIfFileContainsAPrivateRSAKey))
}

func removePubSuffixFromFileName(s string) string {
	return strings.TrimSuffix(s, ".pub")
}

func removePubSuffixFromFileNamesList(files []string) []string {
	return transform(files, removePubSuffixFromFileName)
}

func findKeyPairsBasedOnFileName(privateFiles, publicFiles []string) []string {
	return filter(privateFiles, existsIn(publicFiles))
}

func withoutFileName(targetFileNamesList []string, fileNameToDelete string) []string {
	return filter(targetFileNamesList, not(isEqualTo(fileNameToDelete)))
}

func removeFileNames(targetFileNamesList, fileNameToDelete []string) []string {
	return foldLeft(fileNameToDelete, targetFileNamesList, withoutFileName)
}

func listFilesInHomeSSHDirectory() []string {
	sshDirectory := path.Join(os.Getenv("HOME"), ".ssh")
	return listFilesIn(sshDirectory)
}

func createPublicKeyEntriesFrom(input []string) []KeyEntry {
	return transform(input, func(s string) KeyEntry {
		return createPublicKeyRepresentation(s)
	})
}

func createPrivateKeyEntriesFrom(input []string) []KeyEntry {
	return transform(input, func(s string) KeyEntry {
		return createPrivateKeyRepresentation(s)
	})
}

func privateKeyEntriesFrom(input []string) []KeyEntry {
	return createPrivateKeyEntriesFrom(selectFilesContainingRSAPrivateKeys(input))
}

func publicKeyEntriesFrom(input []string) []KeyEntry {
	return createPublicKeyEntriesFrom(selectFilesContainingRSAPublicKeys(input))
}

type keyEntryPartitioner struct {
	result     []KeyEntry
	publicKeys map[string]*publicKeyRepresentation
}

func (p *keyEntryPartitioner) initializePublicKeyCache(publics []*publicKeyRepresentation) {
	p.publicKeys = map[string]*publicKeyRepresentation{}

	foreach(publics, p.addPublicKeyToCache)
}

func (p *keyEntryPartitioner) publicKeyNameFor(priv *privateKeyRepresentation) string {
	return fmt.Sprintf("%s.pub", priv.path)
}

func (p *keyEntryPartitioner) addResult(r KeyEntry) {
	p.result = append(p.result, r)
}

func (p *keyEntryPartitioner) addPublicKeyResult(r *publicKeyRepresentation) {
	p.addResult(r)
}

func (p *keyEntryPartitioner) potentialPublicKeyFor(priv *privateKeyRepresentation) (pub *publicKeyRepresentation, ok bool) {
	potentialPub, ok := p.publicKeys[p.publicKeyNameFor(priv)]
	return potentialPub, ok
}

func (p *keyEntryPartitioner) addPublicKeyToCache(pub *publicKeyRepresentation) {
	p.publicKeys[pub.path] = pub
}

func (p *keyEntryPartitioner) deleteFromCache(pub *publicKeyRepresentation) {
	delete(p.publicKeys, pub.path)
}

func (p *keyEntryPartitioner) processPrivateKey(priv *privateKeyRepresentation) {
	if potentialPub, ok := p.potentialPublicKeyFor(priv); ok {
		fmt.Println("IS A MATCH!!")
		p.addResult(createKeypairRepresentation(priv, potentialPub))
		p.deleteFromCache(potentialPub)
	} else {
		p.addResult(priv)
	}
}

func (p *keyEntryPartitioner) processPrivateKeys(privates []*privateKeyRepresentation) {
	foreach(privates, p.processPrivateKey)
}

func (p *keyEntryPartitioner) appendRemainingPublicKeys() {
	foreachValue(p.publicKeys, p.addPublicKeyResult)
}

func partitionKeyEntries(privates []*privateKeyRepresentation, publics []*publicKeyRepresentation) []KeyEntry {
	p := &keyEntryPartitioner{}
	p.initializePublicKeyCache(publics)
	p.processPrivateKeys(privates)
	p.appendRemainingPublicKeys()
	return p.result
}
