package ssh

func (k *privateKey) isRSA() bool {
	return k.isAlgorithm(rsaAlgorithm)
}
