package ssh

func (k *privateKey) isEd25519() bool {
	return k.isAlgorithm(ed25519Algorithm)
}
