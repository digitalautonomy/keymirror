package ssh

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/digitalautonomy/keymirror/api"
)

func (s *sshSuite) Test_createPublicKeyRepresentation_createsTheObject() {
	r := createPublicKeyRepresentation("foo_rsa.pub")

	s.Equal("foo_rsa.pub", r.path)

	r = createPublicKeyRepresentation("something else with no pub")

	s.Equal("something else with no pub", r.path)
}

func (s *sshSuite) Test_createPrivateKeyRepresentation_createsTheObject() {
	r := createPrivateKeyRepresentation("foo_rsa")

	s.Equal("foo_rsa", r.path)

	r = createPrivateKeyRepresentation("something else")

	s.Equal("something else", r.path)
}

func (s *sshSuite) Test_createKeypairRepresentation_panicsOnInvalidArguments() {
	s.Panics(func() {
		createKeypairRepresentation(nil, nil)
	}, "panics when both arguments are nil, since this is a developer error")

	s.Panics(func() {
		createKeypairRepresentation(nil, createPublicKeyRepresentation("foo"))
	}, "panics when private key arguments is nil, since this is a developer error")

	s.Panics(func() {
		createKeypairRepresentation(createPrivateKeyRepresentation("something"), nil)
	}, "panics when public key arguments is nil, since this is a developer error")
}

func (s *sshSuite) Test_createKeypairRepresentation_createsAValidRepresentationWhenCorrectArgumentsAreGiven() {
	pub1 := createPublicKeyRepresentation("hello.pub")
	priv1 := createPrivateKeyRepresentation("another name.priv")

	kp1 := createKeypairRepresentation(priv1, pub1)

	s.Equal(priv1, kp1.private)
	s.Equal(pub1, kp1.public)

	pub2 := createPublicKeyRepresentation("another public key.rsa")
	priv2 := createPrivateKeyRepresentation("our secret secret")

	kp2 := createKeypairRepresentation(priv2, pub2)

	s.Equal(priv2, kp2.private)
	s.Equal(pub2, kp2.public)
}

func (s *sshSuite) Test_privateKeyRepresentation_Locations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPrivateKeyRepresentation("")

	s.Empty(pk.Locations())
}

func (s *sshSuite) Test_privateKeyRepresentation_Locations_returnsASliceWithThePathToThePrivateKey() {
	pk := createPrivateKeyRepresentation("/foo/bar/hello.priv")
	s.Equal([]string{"/foo/bar/hello.priv"}, pk.Locations())

	pk = createPrivateKeyRepresentation("another private key")
	s.Equal([]string{"another private key"}, pk.Locations())
}

func (s *sshSuite) Test_privateKeyRepresentation_PublicKeyLocations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPrivateKeyRepresentation("")

	s.Empty(pk.PublicKeyLocations())
}

func (s *sshSuite) Test_privateKeyRepresentation_PublicKeyLocations_returnsAnEmptySliceEvenIfAFilenameIsGiven() {
	pk := createPrivateKeyRepresentation("/foo/bar/hello.priv")
	s.Empty(pk.PublicKeyLocations())

	pk = createPrivateKeyRepresentation("another private key")
	s.Empty(pk.PublicKeyLocations())
}

func (s *sshSuite) Test_privateKeyRepresentation_PrivateKeyLocations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPrivateKeyRepresentation("")

	s.Empty(pk.PrivateKeyLocations())
}

func (s *sshSuite) Test_privateKeyRepresentation_PrivateKeyLocations_returnsThePathOfThePrivateKeyGiven() {
	pk := createPrivateKeyRepresentation("one private key")
	s.Equal([]string{"one private key"}, pk.PrivateKeyLocations())

	pk = createPrivateKeyRepresentation("another private key")
	s.Equal([]string{"another private key"}, pk.PrivateKeyLocations())
}

func (s *sshSuite) Test_privateKeyRepresentation_KeyType_returnsPrivate() {
	pk := createPrivateKeyRepresentation("/foo/bar/hello")

	s.Equal(api.PrivateKeyType, pk.KeyType())
}

func (s *sshSuite) Test_publicKeyRepresentation_Locations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPublicKeyRepresentation("")

	s.Empty(pk.Locations())
}

func (s *sshSuite) Test_publicKeyRepresentation_Locations_returnsASliceWithThePathToThePublicKey() {
	pk := createPublicKeyRepresentation("/foo/bar/hello.pub")
	s.Equal([]string{"/foo/bar/hello.pub"}, pk.Locations())

	pk = createPublicKeyRepresentation("a public RSA key file")
	s.Equal([]string{"a public RSA key file"}, pk.Locations())
}

func (s *sshSuite) Test_publicKeyRepresentation_PrivateKeyLocations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPublicKeyRepresentation("")

	s.Empty(pk.PrivateKeyLocations())
}

func (s *sshSuite) Test_publicKeyRepresentation_PrivateKeyLocations_returnsAnEmptySliceEvenIfAFilenameIsGiven() {
	pk := createPublicKeyRepresentation("/foo/bar/hello.pub")
	s.Empty(pk.PrivateKeyLocations())

	pk = createPublicKeyRepresentation("another public key")
	s.Empty(pk.PrivateKeyLocations())
}

func (s *sshSuite) Test_publicKeyRepresentation_PublicKeyLocations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPublicKeyRepresentation("")

	s.Empty(pk.PublicKeyLocations())
}

func (s *sshSuite) Test_publicKeyRepresentation_PublicKeyLocations_returnsASliceWithThePathToThePublicKey() {
	pk := createPublicKeyRepresentation("/foo/bar/hello.pub")
	s.Equal([]string{"/foo/bar/hello.pub"}, pk.PublicKeyLocations())

	pk = createPublicKeyRepresentation("a public RSA key file")
	s.Equal([]string{"a public RSA key file"}, pk.PublicKeyLocations())
}

func (s *sshSuite) Test_publicKeyRepresentation_KeyType_returnsPublic() {
	pk := createPublicKeyRepresentation("/foo/bar/hello.pub")

	s.Equal(api.PublicKeyType, pk.KeyType())
}

func (s *sshSuite) Test_keypairRepresentation_Locations_returnsAnEmptyList_ifBothPrivateOrPublicKeyHaveEmptyPaths() {
	priv := createPrivateKeyRepresentation("")
	pub := createPublicKeyRepresentation("")

	kp := createKeypairRepresentation(priv, pub)

	s.Empty(kp.Locations())
}

func (s *sshSuite) Test_keypairRepresentation_Locations_returnsOnlyThePublicKeyPath_ifThePrivateKeyPathIsEmpty() {
	priv := createPrivateKeyRepresentation("")
	pub := createPublicKeyRepresentation("a key.pub")

	kp := createKeypairRepresentation(priv, pub)

	s.Equal([]string{"a key.pub"}, kp.Locations())
}

func (s *sshSuite) Test_keypairRepresentation_Locations_returnsOnlyThePrivateKeyPath_ifThePublicKeyPathIsEmpty() {
	priv := createPrivateKeyRepresentation("/something/secret.rsa")
	pub := createPublicKeyRepresentation("")

	kp := createKeypairRepresentation(priv, pub)

	s.Equal([]string{"/something/secret.rsa"}, kp.Locations())
}

func (s *sshSuite) Test_keypairRepresentation_Locations_returnsThePrivateAndPublicPaths() {
	priv1 := createPrivateKeyRepresentation("/home/amnesia/.ssh/foo_rsa")
	pub1 := createPublicKeyRepresentation("/home/amnesia/.ssh/foo_rsa.pub")

	kp1 := createKeypairRepresentation(priv1, pub1)

	s.Equal([]string{"/home/amnesia/.ssh/foo_rsa", "/home/amnesia/.ssh/foo_rsa.pub"}, kp1.Locations())

	priv2 := createPrivateKeyRepresentation("/home/another private key")
	pub2 := createPublicKeyRepresentation("pub.rsa.4096.{{{")

	kp2 := createKeypairRepresentation(priv2, pub2)

	s.Equal([]string{"/home/another private key", "pub.rsa.4096.{{{"}, kp2.Locations())
}

func (s *sshSuite) Test_keypairRepresentation_PrivateKeyLocations_returnsAnEmptyList_ifBothKeysHaveEmptyPaths() {
	priv := createPrivateKeyRepresentation("")
	pub := createPublicKeyRepresentation("")

	kp := createKeypairRepresentation(priv, pub)

	s.Empty(kp.PrivateKeyLocations())
}

func (s *sshSuite) Test_keypairRepresentation_PrivateKeyLocations_returnsAListWithPrivateKeyPaths_ifOnlyThePrivateKeyHasAPath() {
	priv := createPrivateKeyRepresentation("/home/foo/.ssh/privatekey.ed25519")
	pub := createPublicKeyRepresentation("")

	kp := createKeypairRepresentation(priv, pub)

	s.Equal([]string{"/home/foo/.ssh/privatekey.ed25519"}, kp.PrivateKeyLocations())
}

func (s *sshSuite) Test_keypairRepresentation_PrivateKeyLocations_returnsOnlyThePrivatePaths() {
	priv1 := createPrivateKeyRepresentation("/home/amnesia/.ssh/foo_rsa")
	pub1 := createPublicKeyRepresentation("/home/amnesia/.ssh/foo_rsa.pub")

	kp1 := createKeypairRepresentation(priv1, pub1)

	s.Equal([]string{"/home/amnesia/.ssh/foo_rsa"}, kp1.PrivateKeyLocations())

	priv2 := createPrivateKeyRepresentation("/home/another private key")
	pub2 := createPublicKeyRepresentation("pub.rsa.4096.{{{")

	kp2 := createKeypairRepresentation(priv2, pub2)

	s.Equal([]string{"/home/another private key"}, kp2.PrivateKeyLocations())
}

func (s *sshSuite) Test_keypairRepresentation_PublicKeyLocations_returnsAnEmptyList_ifBothKeysHaveEmptyPaths() {
	priv := createPrivateKeyRepresentation("")
	pub := createPublicKeyRepresentation("")

	kp := createKeypairRepresentation(priv, pub)

	s.Empty(kp.PublicKeyLocations())
}

func (s *sshSuite) Test_keypairRepresentation_PublicKeyLocations_returnsAnEmptyList_ifOnlyThePrivateKeyHasAPath() {
	priv := createPrivateKeyRepresentation("/home/foo/.ssh/privatekey.ed25519")
	pub := createPublicKeyRepresentation("")

	kp := createKeypairRepresentation(priv, pub)

	s.Empty(kp.PublicKeyLocations())
}

func (s *sshSuite) Test_keypairRepresentation_PublicKeyLocations_returnsOnlyThePublicPaths() {
	priv1 := createPrivateKeyRepresentation("/home/amnesia/.ssh/foo_rsa")
	pub1 := createPublicKeyRepresentation("/home/amnesia/.ssh/foo_rsa.pub")

	kp1 := createKeypairRepresentation(priv1, pub1)

	s.Equal([]string{"/home/amnesia/.ssh/foo_rsa.pub"}, kp1.PublicKeyLocations())

	priv2 := createPrivateKeyRepresentation("/home/another private key")
	pub2 := createPublicKeyRepresentation("pub.rsa.4096.{{{")

	kp2 := createKeypairRepresentation(priv2, pub2)

	s.Equal([]string{"pub.rsa.4096.{{{"}, kp2.PublicKeyLocations())
}

func (s *sshSuite) Test_pairKeyRepresentation_KeyType_returnsPair() {
	priv := createPrivateKeyRepresentation("/home/another private key")
	pub := createPublicKeyRepresentation("pub.rsa.4096.{{{")

	kp := createKeypairRepresentation(priv, pub)

	s.Equal(api.PairKeyType, kp.KeyType())
}

func decode(b64 string) []byte {
	r, _ := base64.StdEncoding.DecodeString(b64)
	return r
}

const originalKey = "AAAAB3NzaC1yc2EAAAADAQABAAABgQDjXHGj/u27siLdKDii3ijK8xjKwRgYw6f2r+qfk9f5Nc3CD71Hpxndh3kg3GTdiQauPtzJHeFsm2dFW6lxz0nFqIkrxqTV+maJ57bABehbnFP25IxSKSG/8JBqNDdkipQq14mHmCHCmtU4KTzonHqOssvOkjKSs1iPPK0PBONGKXWMrrNm8U9ProLcwtHsNkJNXcq1qhAdbn+ICVcOB19Wibjp5noLZ2NRxxuzJHPwAXuXT22LpVobaeD/HvqoXPXP5R91smlfFvRSgvoRDr+ufeEL8sUO++emu3+ThG+iEvLGeucR1o4+UFweXM9UdPQW4HFrIpvMPHnbg8ysT07vqH/tA0+TxJocFcADxM89+k/4mS7nRZgCnqpe50JuuJYOBWhZxaaanasEY/+lIdni9u3ZTwhPHG8trokPxdpklKQnN2Pg2Y2SX2N8u3s+IYaTO6s9Cs6fokiK475w1UexTMYW0XuJdN8wFRfDMJpZm+ICVc+K6BpxAyABe4u1RPc="
const otherKey = "AAAAB3NzaC1yc2EAAAADAQABAAABgQC6CyfdeOltbKbISAuuvH27pLNxsNsJ18z29jiZLJ5kvJ9kOXXiZxvZW1a394G9YgDpbUFwjbMz4WgkFNPW7+VVdz07JXpxzs4IrCSE8zl944v98k1kwZ5n6jZR3A51jmb55KjvzFeRv2fqbxb7ylV1R5oNpSYu8l7HaWpR0YSjtuZcSahhPZS4hSbgAKrLm+mn6gfLHyYKEeQ0NRpwxybmrM+dEBdR/bs0JxEgJvrfsOahYEbZhL627NQx9F1NfPq2yGr13lLmA7wIIu633WAyxJDGaCfRQTXR67W16Hl+0LIbhe4b5iKJH/7+tf13C2VMewDPCZwDxf6XAwD1XmEC4L64rhGhFp51u9qrbFQnaDAfTRfgm8sS4oxsHMnWc3TRes3PKcQ08CGkSV42EPiQuIyQaUojVcNo+xhCJzqv77/GcwyxCtm3FN/gjlXLxPtFrWjcvd91Z8v/7npjzJDAOUnAV8QFijCzC688oSuoBr5q9G9zcWd7MVY2vpms1L8="

func (s *sshSuite) Test_publicKeyRepresentation_WithDigestContent_returnsTheFingerPrint() {
	key := &publicKeyRepresentation{key: decode(originalKey)}
	fingerprint := key.key
	result := key.WithDigestContent(func(in []byte) []byte {
		return in
	})
	s.Equal(result, fingerprint)

	key = &publicKeyRepresentation{key: decode(originalKey)}
	fingerprint = decode("rPPC5HJW2WNPqThrwNnL7szNtrEC7lUjKLlt0Jnunoo=")
	result = key.WithDigestContent(func(in []byte) []byte {
		res := sha256.Sum256(in)
		return res[:]
	})
	s.Equal(result, fingerprint)

	key = &publicKeyRepresentation{key: decode(otherKey)}
	fingerprint = decode("Az/Dp2M/PXj/fsxRQWWj954BgtKRX8DJ1t7nDrS+TTw=")
	result = key.WithDigestContent(func(in []byte) []byte {
		res := sha256.Sum256(in)
		return res[:]
	})
	s.Equal(result, fingerprint)
}
