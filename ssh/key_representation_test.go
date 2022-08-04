package ssh

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/digitalautonomy/keymirror/api"
)

func (s *sshSuite) Test_createPublicKeyRepresentation_createsTheObject() {
	r := createPublicKeyRepresentationForTest("foo_rsa.pub", "")

	s.Equal("foo_rsa.pub", r.path)

	r = createPublicKeyRepresentationForTest("something else with no pub", "")

	s.Equal("something else with no pub", r.path)
}

func (s *sshSuite) Test_createPrivateKeyRepresentation_createsTheObject() {
	r := createPrivateKeyRepresentationForTest("foo_rsa")

	s.Equal("foo_rsa", r.path)

	r = createPrivateKeyRepresentationForTest("something else")

	s.Equal("something else", r.path)
}

func (s *sshSuite) Test_createKeypairRepresentation_panicsOnInvalidArguments() {
	s.Panics(func() {
		createKeypairRepresentation(nil, nil)
	}, "panics when both arguments are nil, since this is a developer error")

	s.Panics(func() {
		createKeypairRepresentation(nil, createPublicKeyRepresentationForTest("foo", "ForTest"))
	}, "panics when private key arguments is nil, since this is a developer error")

	s.Panics(func() {
		createKeypairRepresentation(createPrivateKeyRepresentationForTest("something"), nil)
	}, "panics when public key arguments is nil, since this is a developer error")
}

func (s *sshSuite) Test_createKeypairRepresentation_createsAValidRepresentationWhenCorrectArgumentsAreGiven() {
	pub1 := createPublicKeyRepresentationForTest("hello.pub", "ForTest")
	priv1 := createPrivateKeyRepresentationForTest("another name.priv")

	kp1 := createKeypairRepresentation(priv1, pub1)

	s.Equal(priv1, kp1.private)
	s.Equal(pub1, kp1.public)

	pub2 := createPublicKeyRepresentationForTest("another public key.rsa", "")
	priv2 := createPrivateKeyRepresentationForTest("our secret secret")

	kp2 := createKeypairRepresentation(priv2, pub2)

	s.Equal(priv2, kp2.private)
	s.Equal(pub2, kp2.public)
}

func (s *sshSuite) Test_privateKeyRepresentation_Locations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPrivateKeyRepresentationForTest("")

	s.Empty(pk.Locations())
}

func (s *sshSuite) Test_privateKeyRepresentation_Locations_returnsASliceWithThePathToThePrivateKey() {
	pk := createPrivateKeyRepresentationForTest("/foo/bar/hello.priv")
	s.Equal([]string{"/foo/bar/hello.priv"}, pk.Locations())

	pk = createPrivateKeyRepresentationForTest("another private key")
	s.Equal([]string{"another private key"}, pk.Locations())
}

func (s *sshSuite) Test_privateKeyRepresentation_PublicKeyLocations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPrivateKeyRepresentationForTest("")

	s.Empty(pk.PublicKeyLocations())
}

func (s *sshSuite) Test_privateKeyRepresentation_PublicKeyLocations_returnsAnEmptySliceEvenIfAFilenameIsGiven() {
	pk := createPrivateKeyRepresentationForTest("/foo/bar/hello.priv")
	s.Empty(pk.PublicKeyLocations())

	pk = createPrivateKeyRepresentationForTest("another private key")
	s.Empty(pk.PublicKeyLocations())
}

func (s *sshSuite) Test_privateKeyRepresentation_PrivateKeyLocations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPrivateKeyRepresentationForTest("")

	s.Empty(pk.PrivateKeyLocations())
}

func (s *sshSuite) Test_privateKeyRepresentation_PrivateKeyLocations_returnsThePathOfThePrivateKeyGiven() {
	pk := createPrivateKeyRepresentationForTest("one private key")
	s.Equal([]string{"one private key"}, pk.PrivateKeyLocations())

	pk = createPrivateKeyRepresentationForTest("another private key")
	s.Equal([]string{"another private key"}, pk.PrivateKeyLocations())
}

func (s *sshSuite) Test_privateKeyRepresentation_KeyType_returnsPrivate() {
	pk := createPrivateKeyRepresentationForTest("/foo/bar/hello")

	s.Equal(api.PrivateKeyType, pk.KeyType())
}

func (s *sshSuite) Test_publicKeyRepresentation_Locations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPublicKeyRepresentationForTest("", "")

	s.Empty(pk.Locations())
}

func (s *sshSuite) Test_publicKeyRepresentation_Locations_returnsASliceWithThePathToThePublicKey() {
	pk := createPublicKeyRepresentationForTest("/foo/bar/hello.pub", "")
	s.Equal([]string{"/foo/bar/hello.pub"}, pk.Locations())

	pk = createPublicKeyRepresentationForTest("a public RSA key file", "")
	s.Equal([]string{"a public RSA key file"}, pk.Locations())
}

func (s *sshSuite) Test_publicKeyRepresentation_PrivateKeyLocations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPublicKeyRepresentationForTest("", "")

	s.Empty(pk.PrivateKeyLocations())
}

func (s *sshSuite) Test_publicKeyRepresentation_PrivateKeyLocations_returnsAnEmptySliceEvenIfAFilenameIsGiven() {
	pk := createPublicKeyRepresentationForTest("/foo/bar/hello.pub", "")
	s.Empty(pk.PrivateKeyLocations())

	pk = createPublicKeyRepresentationForTest("another public key", "")
	s.Empty(pk.PrivateKeyLocations())
}

func (s *sshSuite) Test_publicKeyRepresentation_PublicKeyLocations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPublicKeyRepresentationForTest("", "")

	s.Empty(pk.PublicKeyLocations())
}

func (s *sshSuite) Test_publicKeyRepresentation_PublicKeyLocations_returnsASliceWithThePathToThePublicKey() {
	pk := createPublicKeyRepresentationForTest("/foo/bar/hello.pub", "")
	s.Equal([]string{"/foo/bar/hello.pub"}, pk.PublicKeyLocations())

	pk = createPublicKeyRepresentationForTest("a public RSA key file", "")
	s.Equal([]string{"a public RSA key file"}, pk.PublicKeyLocations())
}

func (s *sshSuite) Test_publicKeyRepresentation_KeyType_returnsPublic() {
	pk := createPublicKeyRepresentationForTest("/foo/bar/hello.pub", "")

	s.Equal(api.PublicKeyType, pk.KeyType())
}

func (s *sshSuite) Test_keypairRepresentation_Locations_returnsAnEmptyList_ifBothPrivateOrPublicKeyHaveEmptyPaths() {
	priv := createPrivateKeyRepresentationForTest("")
	pub := createPublicKeyRepresentationForTest("", "")

	kp := createKeypairRepresentation(priv, pub)

	s.Empty(kp.Locations())
}

func (s *sshSuite) Test_keypairRepresentation_Locations_returnsOnlyThePublicKeyPath_ifThePrivateKeyPathIsEmpty() {
	priv := createPrivateKeyRepresentationForTest("")
	pub := createPublicKeyRepresentationForTest("a key.pub", "")

	kp := createKeypairRepresentation(priv, pub)

	s.Equal([]string{"a key.pub"}, kp.Locations())
}

func (s *sshSuite) Test_keypairRepresentation_Locations_returnsOnlyThePrivateKeyPath_ifThePublicKeyPathIsEmpty() {
	priv := createPrivateKeyRepresentationForTest("/something/secret.rsa")
	pub := createPublicKeyRepresentationForTest("", "")

	kp := createKeypairRepresentation(priv, pub)

	s.Equal([]string{"/something/secret.rsa"}, kp.Locations())
}

func (s *sshSuite) Test_keypairRepresentation_Locations_returnsThePrivateAndPublicPaths() {
	priv1 := createPrivateKeyRepresentationForTest("/home/amnesia/.ssh/foo_rsa")
	pub1 := createPublicKeyRepresentationForTest("/home/amnesia/.ssh/foo_rsa.pub", "")

	kp1 := createKeypairRepresentation(priv1, pub1)

	s.Equal([]string{"/home/amnesia/.ssh/foo_rsa", "/home/amnesia/.ssh/foo_rsa.pub"}, kp1.Locations())

	priv2 := createPrivateKeyRepresentationForTest("/home/another private key")
	pub2 := createPublicKeyRepresentationForTest("pub.rsa.4096.{{{", "")

	kp2 := createKeypairRepresentation(priv2, pub2)

	s.Equal([]string{"/home/another private key", "pub.rsa.4096.{{{"}, kp2.Locations())
}

func (s *sshSuite) Test_keypairRepresentation_PrivateKeyLocations_returnsAnEmptyList_ifBothKeysHaveEmptyPaths() {
	priv := createPrivateKeyRepresentationForTest("")
	pub := createPublicKeyRepresentationForTest("", "")

	kp := createKeypairRepresentation(priv, pub)

	s.Empty(kp.PrivateKeyLocations())
}

func (s *sshSuite) Test_keypairRepresentation_PrivateKeyLocations_returnsAListWithPrivateKeyPaths_ifOnlyThePrivateKeyHasAPath() {
	priv := createPrivateKeyRepresentationForTest("/home/foo/.ssh/privatekey.ed25519")
	pub := createPublicKeyRepresentationForTest("", "")

	kp := createKeypairRepresentation(priv, pub)

	s.Equal([]string{"/home/foo/.ssh/privatekey.ed25519"}, kp.PrivateKeyLocations())
}

func (s *sshSuite) Test_keypairRepresentation_PrivateKeyLocations_returnsOnlyThePrivatePaths() {
	priv1 := createPrivateKeyRepresentationForTest("/home/amnesia/.ssh/foo_rsa")
	pub1 := createPublicKeyRepresentationForTest("/home/amnesia/.ssh/foo_rsa.pub", "")

	kp1 := createKeypairRepresentation(priv1, pub1)

	s.Equal([]string{"/home/amnesia/.ssh/foo_rsa"}, kp1.PrivateKeyLocations())

	priv2 := createPrivateKeyRepresentationForTest("/home/another private key")
	pub2 := createPublicKeyRepresentationForTest("pub.rsa.4096.{{{", "")

	kp2 := createKeypairRepresentation(priv2, pub2)

	s.Equal([]string{"/home/another private key"}, kp2.PrivateKeyLocations())
}

func (s *sshSuite) Test_keypairRepresentation_PublicKeyLocations_returnsAnEmptyList_ifBothKeysHaveEmptyPaths() {
	priv := createPrivateKeyRepresentationForTest("")
	pub := createPublicKeyRepresentationForTest("", "")

	kp := createKeypairRepresentation(priv, pub)

	s.Empty(kp.PublicKeyLocations())
}

func (s *sshSuite) Test_keypairRepresentation_PublicKeyLocations_returnsAnEmptyList_ifOnlyThePrivateKeyHasAPath() {
	priv := createPrivateKeyRepresentationForTest("/home/foo/.ssh/privatekey.ed25519")
	pub := createPublicKeyRepresentationForTest("", "")

	kp := createKeypairRepresentation(priv, pub)

	s.Empty(kp.PublicKeyLocations())
}

func (s *sshSuite) Test_keypairRepresentation_PublicKeyLocations_returnsOnlyThePublicPaths() {
	priv1 := createPrivateKeyRepresentationForTest("/home/amnesia/.ssh/foo_rsa")
	pub1 := createPublicKeyRepresentationForTest("/home/amnesia/.ssh/foo_rsa.pub", "")

	kp1 := createKeypairRepresentation(priv1, pub1)

	s.Equal([]string{"/home/amnesia/.ssh/foo_rsa.pub"}, kp1.PublicKeyLocations())

	priv2 := createPrivateKeyRepresentationForTest("/home/another private key")
	pub2 := createPublicKeyRepresentationForTest("pub.rsa.4096.{{{", "")

	kp2 := createKeypairRepresentation(priv2, pub2)

	s.Equal([]string{"pub.rsa.4096.{{{"}, kp2.PublicKeyLocations())
}

func (s *sshSuite) Test_pairKeyRepresentation_KeyType_returnsPair() {
	priv := createPrivateKeyRepresentationForTest("/home/another private key")
	pub := createPublicKeyRepresentationForTest("pub.rsa.4096.{{{", "")

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
	result := key.WithDigestContent(identity[[]byte])
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

func identity[T any](v T) T {
	return v
}

func (s *sshSuite) Test_keypairRepresentation_WithDigestContent_returnsTheFingerPrint() {
	keyPair := &keypairRepresentation{
		public: &publicKeyRepresentation{key: decode(originalKey)},
	}

	s.Equal(
		keyPair.WithDigestContent(identity[[]byte]),
		keyPair.public.WithDigestContent(identity[[]byte]))

	keyPair = &keypairRepresentation{
		public: &publicKeyRepresentation{key: decode(otherKey)},
	}

	s.Equal(
		keyPair.WithDigestContent(identity[[]byte]),
		keyPair.public.WithDigestContent(identity[[]byte]))
}
