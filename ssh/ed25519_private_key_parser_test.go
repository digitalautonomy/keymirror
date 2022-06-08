package ssh

const correctEd25519PrivateKey = `
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACC7eNn1eQ/DPPtfZUAie2p9I1TAuj91YucOlbHyxV7hygAAAJC+YoKmvmKC
pgAAAAtzc2gtZWQyNTUxOQAAACC7eNn1eQ/DPPtfZUAie2p9I1TAuj91YucOlbHyxV7hyg
AAAEBVI12MKVSate/Pvx/nqIe2B4/J3Y8qURPhFGcUZyEtgbt42fV5D8M8+19lQCJ7an0j
VMC6P3Vi5w6VsfLFXuHKAAAACmZhdXN0b0BDQUQBAgM=
-----END OPENSSH PRIVATE KEY-----
`

func (s *sshSuite) Test_parsePrivateKey_AStringContainingACorrectEd25519KeyShouldBeConsideredAPrivateKey() {
	pk := correctEd25519PrivateKey

	a, _ := accessWithTestLogging()
	priv, ok := a.parsePrivateKey(pk)

	s.True(ok)
	s.Equal("ssh-ed25519", priv.algorithm)
}
