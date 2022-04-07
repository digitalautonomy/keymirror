package ssh

import (
	"fmt"
	"github.com/digitalautonomy/keymirror/api"
)

type keyEntryPartitioner struct {
	result     []api.KeyEntry
	publicKeys map[string]*publicKeyRepresentation
}

func (p *keyEntryPartitioner) initializePublicKeyCache(publics []*publicKeyRepresentation) {
	p.publicKeys = map[string]*publicKeyRepresentation{}

	foreach(publics, p.addPublicKeyToCache)
}

func (p *keyEntryPartitioner) publicKeyNameFor(priv *privateKeyRepresentation) string {
	return fmt.Sprintf("%s.pub", priv.path)
}

func (p *keyEntryPartitioner) addResult(r api.KeyEntry) {
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

func partitionKeyEntries(privates []*privateKeyRepresentation, publics []*publicKeyRepresentation) []api.KeyEntry {
	p := &keyEntryPartitioner{}
	p.initializePublicKeyCache(publics)
	p.processPrivateKeys(privates)
	p.appendRemainingPublicKeys()
	return p.result
}
