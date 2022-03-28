package ssh

type privateKeyRepresentation struct {
	path string
}

type publicKeyRepresentation struct {
	path string
}

type keypairRepresentation struct {
	private *privateKeyRepresentation
	public  *publicKeyRepresentation
}

func createPublicKeyRepresentation(path string) *publicKeyRepresentation {
	return &publicKeyRepresentation{
		path: path,
	}
}

func createPrivateKeyRepresentation(path string) *privateKeyRepresentation {
	return &privateKeyRepresentation{
		path: path,
	}
}

// createKeypairRepresentation creates a keypair from the given public and private keys
// it is NOT acceptable to send in nil as any of the arguments - this is a developer error
// and will result in a panic
func createKeypairRepresentation(private *privateKeyRepresentation, public *publicKeyRepresentation) *keypairRepresentation {
	if private == nil {
		panic("private key representation argument is nil, which is a developer error")
	}
	if public == nil {
		panic("public key representation argument is nil, which is a developer error")
	}
	return &keypairRepresentation{
		private: private,
		public:  public,
	}
}

// Locations implement the KeyEntry interface
func (k *privateKeyRepresentation) Locations() []string {
	if k.path == "" {
		return nil
	}
	return []string{k.path}
}

// Locations implement the KeyEntry interface
func (k *publicKeyRepresentation) Locations() []string {
	if k.path == "" {
		return nil
	}
	return []string{k.path}
}

// Locations implement the KeyEntry interface
func (k *keypairRepresentation) Locations() []string {
	return append(k.private.Locations(), k.public.Locations()...)
}
