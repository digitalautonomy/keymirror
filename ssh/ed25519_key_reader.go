package ssh

const ed25519Algorithm = "ssh-ed25519"

func (a *access) isEd25519PrivateKey(pk string) bool {
	priv, ok := a.parsePrivateKey(pk)
	if !ok {
		return false
	}
	return priv.isEd25519()
}

func (k *publicKey) isEd25519() bool {
	return k.isAlgorithm(ed25519Algorithm)
}

func isEd25519PublicKey(k string) bool {
	pub, ok := parsePublicKey(k)
	if !ok {
		return false
	}
	return pub.isEd25519()
}
