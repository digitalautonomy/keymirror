package ssh

const rsaAlgorithm = "ssh-rsa"

func (k *publicKey) isRSA() bool {
	return k.isAlgorithm(rsaAlgorithm)
}

func isRSAPublicKey(k string) bool {
	pub, ok := parsePublicKey(k)
	if !ok {
		return false
	}
	return pub.isRSA()
}

func (a *access) isRSAPrivateKey(pk string) bool {
	priv, ok := a.parsePrivateKey(pk)
	if !ok {
		return false
	}
	return priv.isRSA()
}
