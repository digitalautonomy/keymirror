package ssh

import (
	"github.com/digitalautonomy/keymirror/api"
)

type privateKeyRepresentation struct {
	path              string
	passwordProtected bool
	size              int
	algorithm         api.Algorithm
}

type publicKeyRepresentation struct {
	path      string
	key       []byte
	size      int
	algorithm api.Algorithm
	userID    string
}

type keypairRepresentation struct {
	private *privateKeyRepresentation
	public  *publicKeyRepresentation
}

func createPublicKeyRepresentationFromPublicKey(key *publicKey) *publicKeyRepresentation {
	return &publicKeyRepresentation{
		path:      key.location,
		key:       key.key,
		size:      key.size,
		algorithm: translateSshAlgorithmToExternalAlgorithm(key.algorithm),
		userID:    key.comment,
	}
}

var algorithms = map[string]api.Algorithm{
	rsaAlgorithm:     api.RSA,
	ed25519Algorithm: api.Ed25519,
}

func translateSshAlgorithmToExternalAlgorithm(algo string) api.Algorithm {
	return algorithms[algo]
}

func createPrivateKeyRepresentationFromPrivateKey(key *privateKey) *privateKeyRepresentation {
	return &privateKeyRepresentation{
		path:              key.path,
		passwordProtected: key.passwordProtected,
		size:              key.size,
		algorithm:         translateSshAlgorithmToExternalAlgorithm(key.algorithm),
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

func nilOrStringSlice(s string) []string {
	if s == "" {
		return nil
	}
	return []string{s}
}

// Locations implement the KeyEntry interface
func (k *privateKeyRepresentation) Locations() []string {
	return nilOrStringSlice(k.path)
}

func (k *privateKeyRepresentation) PrivateKeyLocations() []string {
	return k.Locations()
}

func (k *privateKeyRepresentation) KeyType() api.KeyType {
	return api.PrivateKeyType
}

// PublicKeyLocations implement the KeyEntry interface
func (k *privateKeyRepresentation) PublicKeyLocations() []string {
	return nil
}

func (k *privateKeyRepresentation) IsPasswordProtected() bool {
	return k.passwordProtected
}

func (k *privateKeyRepresentation) Size() int {
	return k.size
}

func (k *privateKeyRepresentation) Algorithm() api.Algorithm {
	return k.algorithm
}

// Locations implement the KeyEntry interface
func (k *publicKeyRepresentation) Locations() []string {
	return nilOrStringSlice(k.path)
}

func (k *publicKeyRepresentation) PrivateKeyLocations() []string {
	return nil
}

// PublicKeyLocations implement the KeyEntry interface
func (k *publicKeyRepresentation) PublicKeyLocations() []string {
	return k.Locations()
}

func (k *publicKeyRepresentation) KeyType() api.KeyType {
	return api.PublicKeyType
}

func (k *publicKeyRepresentation) WithDigestContent(f func([]byte) []byte) []byte {
	return f(k.key)
}

func (k *publicKeyRepresentation) Size() int {
	return k.size
}

func (k *publicKeyRepresentation) Algorithm() api.Algorithm {
	return k.algorithm
}

func (k *publicKeyRepresentation) UserID() string {
	return k.userID
}

// Locations implement the KeyEntry interface
func (k *keypairRepresentation) Locations() []string {
	return append(k.private.Locations(), k.public.Locations()...)
}

func (k *keypairRepresentation) PrivateKeyLocations() []string {
	return k.private.PrivateKeyLocations()
}

// PublicKeyLocations implement the KeyEntry interface
func (k *keypairRepresentation) PublicKeyLocations() []string {
	return k.public.PublicKeyLocations()
}

func (k *keypairRepresentation) KeyType() api.KeyType {
	return api.PairKeyType
}

func (k *keypairRepresentation) WithDigestContent(f func([]byte) []byte) []byte {
	return k.public.WithDigestContent(f)
}

func (k *keypairRepresentation) IsPasswordProtected() bool {
	return k.private.passwordProtected
}

func (k *keypairRepresentation) Size() int {
	return k.public.size
}

func (k *keypairRepresentation) Algorithm() api.Algorithm {
	return k.public.Algorithm()
}

func (k *keypairRepresentation) UserID() string {
	return k.public.userID
}
