package ssh

import "github.com/digitalautonomy/keymirror/api"

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
