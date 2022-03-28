package ssh

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

func (s *sshSuite) Test_publicKeyRepresentation_Locations_returnsAnEmptySliceIfNoFilenameIsGiven() {
	pk := createPublicKeyRepresentation("")

	s.Empty(pk.Locations())
}

func (s *sshSuite) Test_publicKeyRepresentation_Locations_returnsASliceWithThePathToThePrivateKey() {
	pk := createPublicKeyRepresentation("/foo/bar/hello.pub")
	s.Equal([]string{"/foo/bar/hello.pub"}, pk.Locations())

	pk = createPublicKeyRepresentation("a public RSA key file")
	s.Equal([]string{"a public RSA key file"}, pk.Locations())
}

func (s *sshSuite) Test_keypairRepresentation_Locations_returnsAnEmptyList_ifNeitherPrivateOrPublicKeyHasEmptyPath() {
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
