package ssh

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type sshSuite struct {
	suite.Suite
	tdir string
}

func TestSSHSuite(t *testing.T) {
	suite.Run(t, new(sshSuite))
}
func (s *sshSuite) SetupTest() {
	s.tdir = s.T().TempDir()
}

func (s *sshSuite) Test_ListsAllTheFilesInSpecifiedDirectory() {
	r := rand.Int()

	expected := []string{"id_rsa.pub", fmt.Sprintf("id_rsa%d", r)}

	for _, f := range expected {
		s.createFileWithContent(s.tdir, f, "some content")
	}

	s.Equal(expected, listFilesIn(s.tdir))
}

func (s *sshSuite) Test_ListsNoFilesInADirectoryThatDoesntExist() {
	s.Equal([]string{}, listFilesIn("directory that hopefully doesnt exist"))
}

func (s *sshSuite) Test_ParseAStringAsAnSSHPublicKeyRepresentation() {
	k := ""
	_, ok := parsePublicKey(k)
	s.Require().False(ok, "An empty string is not a valid SSH public key representation")

	k = "ssh-rsa dmFsaWQgYmFzZTY0IHN0cmluZw== batman@debian"
	pub, ok := parsePublicKey(k)
	s.Require().True(ok, "Should parse a valid SSH RSA public key representation")
	s.Equal("ssh-rsa", pub.algorithm)

	k = "ssh-ecdsa bla2 robin@debian"
	pub, ok = parsePublicKey(k)
	s.Require().True(ok, "Should parse a valid SSH public key representation with a different key type")
	s.Equal("ssh-ecdsa", pub.algorithm)

	k = "ssh-rsa"
	_, ok = parsePublicKey(k)
	s.Require().False(ok, "A string with only one field is not a valid SSH public key representation")

	k = "ssh-rsa  "
	_, ok = parsePublicKey(k)
	s.False(ok, "Since more than one whitespace character serve as one single separator, this example only has one column, and is thus not valid")

	k = "ssh-rsa  b3RoZXIgdmFsaWQgc3RyaW5n foo@debian"
	_, ok = parsePublicKey(k)
	s.True(ok, "More than one whitespace character serves as one single separator between columns")

	k = "ssh-rsa\tb3RoZXIgdmFsaWQgc3RyaW5n foo@debian"
	_, ok = parsePublicKey(k)
	s.True(ok, "A tab can be a separator for columns")

	k = "ssh-rsa   \t  \t  \t b3RoZXIgdmFsaWQgc3RyaW5n foo@debian"
	pub, ok = parsePublicKey(k)
	s.True(ok, "A mix of tabs and spaces serve as one single separator")
	s.Equal(decode("b3RoZXIgdmFsaWQgc3RyaW5n"), pub.key)

	k = "ssh-rsa AAQQ foo@debian foo2@debian"
	pub, ok = parsePublicKey(k)
	s.True(ok, "More than one comment is acceptable in an SSH public key")
	s.Equal("foo@debian foo2@debian", pub.comment)

	k = "ssh-rsa AAQQ"
	_, ok = parsePublicKey(k)
	s.True(ok, "An SSH public key without a comment is still acceptable")
}

func (s *sshSuite) Test_parsePublicKey_doesNotParseAKeyThatIsNotAValidBase64() {
	k := "ssh-rsa b batman@debian"
	_, ok := parsePublicKey(k)
	s.Require().False(ok, "Should not accept a non-base64 key")
}

func (s *sshSuite) Test_parsePublicKey_parsesABase64Key() {
	k := "ssh-rsa YSB2YWxpZCBiYXNlNjQga2V5 batman@debian"
	key, ok := parsePublicKey(k)
	s.Require().True(ok, "Should accept a Base64 key")
	s.Require().Equal([]byte("a valid base64 key"), key.key)
}

func (s *sshSuite) Test_CheckIfThePublicKeyTypeIdentifierIsRSA() {
	pub := publicKey{}
	s.False(pub.isRSA(), "An empty key is not an RSA key")

	pub = publicKey{algorithm: rsaAlgorithm}
	s.True(pub.isRSA(), "A key with the algorithm identifier ssh-rsa is an RSA key")

	pub = publicKey{algorithm: "ssh-ecdsa"}
	s.False(pub.isRSA(), "A key with the algorithm identifier ssh-ecdsa is not an RSA key")
}

func (s *sshSuite) Test_isRSAPublicKey_checkIfAStringHasTheFormatOfAnRSAPublicKey() {
	k := ""
	s.False(isRSAPublicKey(k), "An empty string is not an SSH public key representation thus it is not an RSA key")

	k = "ssh-rsa AAQQ"
	s.True(isRSAPublicKey(k), "A string with the algorithm identifier ssh-rsa is an RSA key")

	k = "ssh-ecdsa AAQQ"
	s.False(isRSAPublicKey(k), "A string with the algorithm identifier ssh-ecdsa is not an RSA key")
}

func (s *sshSuite) Test_parsePrivateKey_AnEmptyStringIsNotAPrivateKey() {
	pk := ""

	a, _ := accessWithTestLogging()
	_, ok := a.parsePrivateKey(pk)

	s.False(ok)
}

func (s *sshSuite) Test_parsePrivateKey_AStringThatDoesNotFollowThePEMFormatShouldNotBeConsideredAPrivateKey() {
	pk := "This is not a RSA Private Key and this does not follow the PEM format"

	a, _ := accessWithTestLogging()
	_, ok := a.parsePrivateKey(pk)

	s.False(ok)
}

func (s *sshSuite) Test_parsePrivateKey_AStringThatRepresentsAPublicKeyInThePEMFormatShouldNotBeConsideredAPrivateKey() {
	pk := `
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x
qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY
nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q
TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/
a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9
PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ
yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUk
AIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ==
-----END PUBLIC KEY-----
`

	a, _ := accessWithTestLogging()
	_, ok := a.parsePrivateKey(pk)

	s.False(ok)
}

func (s *sshSuite) Test_parsePrivateKey_AStringThatRepresentsAnSSLPrivateInThePEMFormatShouldNotBeConsideredAnSSHPrivateKey() {
	pk := `
-----BEGIN OLAS FANCY OPENSSH PRIVATE KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x
qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY
nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q
TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/
a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9
PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ
yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUk
AIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ==
-----END OLAS FANCY OPENSSH PRIVATE KEY-----
`

	a, _ := accessWithTestLogging()
	_, ok := a.parsePrivateKey(pk)

	s.False(ok)
}

func (s *sshSuite) Test_parsePrivateKey_AStringContainingMismatchedPEMTypeShouldNotBeConsideredAPrivateKey() {
	pk := `
-----BEGIN OPENSSH PRIVATE KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x
qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY
nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q
TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/
a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9
PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ
yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUk
AIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ==
-----END PUBLIC KEY-----
`

	a, _ := accessWithTestLogging()
	_, ok := a.parsePrivateKey(pk)

	s.False(ok)
}

func (s *sshSuite) Test_parsePrivateKey_AStringContainingCorrectPEMFormatButMalformedBase64ContentShouldNotBeConsideredAPrivateKey() {
	pk := `
-----BEGIN OPENSSH PRIVATE KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x
qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY
nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q
TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/
a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9
PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ
yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUk
AIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ=
-----END OPENSSH PRIVATE KEY-----
`

	a, _ := accessWithTestLogging()
	_, ok := a.parsePrivateKey(pk)

	s.False(ok)
}

const sshMagicValueEncoded = "b3BlbnNzaC1rZXktdjE"

const correctECDSASSHPrivateKey = `
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAaAAAABNlY2RzYS
1zaGEyLW5pc3RwMjU2AAAACG5pc3RwMjU2AAAAQQR9WZPeBSvixkhjQOh9yCXXlEx5CN9M
yh94CJJ1rigf8693gc90HmahIR5oMGHwlqMoS7kKrRw+4KpxqsF7LGvxAAAAqJZtgRuWbY
EbAAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH1Zk94FK+LGSGNA
6H3IJdeUTHkI30zKH3gIknWuKB/zr3eBz3QeZqEhHmgwYfCWoyhLuQqtHD7gqnGqwXssa/
EAAAAgBzKpRmMyXZ4jnSt3ARz0ul6R79AXAr5gQqDAmoFeEKwAAAAOYWpAYm93aWUubG9j
YWwBAg==
-----END OPENSSH PRIVATE KEY-----
`

func (s *sshSuite) Test_parsePrivateKey_AStringContainingAWellFormedOpenSSHPrivateKeyShouldBeConsideredAPrivateKey() {
	pk := correctECDSASSHPrivateKey

	a, _ := accessWithTestLogging()
	priv, ok := a.parsePrivateKey(pk)

	s.True(ok)
	s.Equal("ecdsa-sha2-nistp256", priv.algorithm)
}

func corruptMagicValue(s string) string {
	return strings.ReplaceAll(s, sshMagicValueEncoded, "b3BlbnLzaC1rZXktdjE")
}

func (s *sshSuite) Test_parsePrivateKey_AStringContainingABinaryWithoutACorrectSSHMagicValueIsNotAValidPrivateKey() {
	pk := corruptMagicValue(correctECDSASSHPrivateKey)

	a, _ := accessWithTestLogging()
	_, ok := a.parsePrivateKey(pk)

	s.False(ok)
}

func (s *sshSuite) Test_parsePrivateKey_AStringContainingAVeryShortBase64StringIsNotAValidPrivateKey() {
	pk := `
-----BEGIN OPENSSH PRIVATE KEY-----
b3Bl
-----END OPENSSH PRIVATE KEY-----
`
	a, _ := accessWithTestLogging()
	_, ok := a.parsePrivateKey(pk)

	s.False(ok)
}

const correctRSASSHPrivateKey = `
-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----
`

const correctRSASSHPrivateKeyOther = `
-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----
`

func (s *sshSuite) Test_parsePrivateKey_AStringContainingAWellFormedRSAOpenSSHPrivateKeyShouldBeConsideredAPrivateKey() {
	pk := correctRSASSHPrivateKey

	a, _ := accessWithTestLogging()
	priv, ok := a.parsePrivateKey(pk)

	s.True(ok)
	s.Equal("ssh-rsa", priv.algorithm)
}

func (s *sshSuite) Test_read32BitNumber_AnEmptyByteSliceIsNotAValid32BitNumber() {
	input := []byte{}

	_, _, ok := read32BitNumber(input)

	s.False(ok)
}

func (s *sshSuite) Test_read32BitNumber_AnInputWithLessThan32BitsIsNotAValid32BitNumber() {
	input := []byte{8, 12, 13}

	_, _, ok := read32BitNumber(input)

	s.False(ok)
}

func (s *sshSuite) Test_read32BitNumber_AnInputWithFourZerosIsAValid32BitNumber() {
	input := []byte{0, 0, 0, 0}

	_, _, ok := read32BitNumber(input)

	s.True(ok)
}

func (s *sshSuite) Test_read32BitNumber_InputWithNonZeroValueReturnsValue() {
	input := []byte{13, 8, 16, 12}

	v, _, ok := read32BitNumber(input)

	s.True(ok)
	s.Equal(uint32(218632204), v)
}

func (s *sshSuite) Test_read32BitNumber_InputWithAnotherNonZeroValueReturnsThatValue() {
	input := []byte{0, 0, 0, 1}

	v, _, ok := read32BitNumber(input)

	s.True(ok)
	s.Equal(uint32(1), v)
}

func (s *sshSuite) Test_read32BitNumber_WhenOkWeHaveARest() {
	input := []byte{0, 0, 0, 1}

	_, r, ok := read32BitNumber(input)

	s.True(ok)
	s.Empty(r)
}

func (s *sshSuite) Test_read32BitNumber_WithALongerInputTheRestShouldBeTheRemainingBytes() {
	input := []byte{0, 0, 0, 1, 1, 2, 3}
	_, r, ok := read32BitNumber(input)
	s.True(ok)
	s.Equal([]byte{1, 2, 3}, r)

	input = []byte{0, 0, 0, 1, 4, 5, 6}
	_, r, ok = read32BitNumber(input)
	s.True(ok)
	s.Equal([]byte{4, 5, 6}, r)
}

func (s *sshSuite) Test_readLengthBytes_AnEmptyInputHasNotAValidLength() {
	input := []byte{}
	_, _, ok := readLengthBytes(input)
	s.False(ok)
}

func (s *sshSuite) Test_readLengthBytes_AnInputWithoutSufficientRestReturnsFalse() {
	input := []byte{0, 0, 0, 12, 8}
	_, _, ok := readLengthBytes(input)
	s.False(ok)
}

func (s *sshSuite) Test_readLengthBytes_AnInputWithEnoughRestReturnsAtLeastTheValue() {
	input := []byte{0, 0, 0, 1, 8}
	v, _, ok := readLengthBytes(input)
	s.True(ok)
	s.Equal([]byte{8}, v)
}

func (s *sshSuite) Test_readLengthBytes_AnInputWithMoreThanEnoughRestReturnsTheValueAndARest() {
	input := []byte{0, 0, 0, 1, 8, 12, 8, 15}
	v, r, ok := readLengthBytes(input)
	s.True(ok)
	s.Equal([]byte{8}, v)
	s.Equal([]byte{12, 8, 15}, r)
}

func (s *sshSuite) Test_extractPrivateKeyAlgorithm_FromAnEmptyInputItIsNotPossibleToExtractAValidPrivateKeyAlgorithm() {
	input := []byte{}
	_, ok := extractKeyAlgorithm(input)
	s.False(ok)
}

func (s *sshSuite) Test_extractPrivateKeyAlgorithm_ANotLongEnoughByteSliceShouldReturnNotOk() {
	input := []byte{0, 0, 0, 7}
	input = append(input, []byte("ssh-r")...)

	_, ok := extractKeyAlgorithm(input)
	s.False(ok)
}

func (s *sshSuite) Test_extractPrivateKeyAlgorithm_ALongEnoughByteSliceShouldReturnAValidPrivateKeyAlgorithm() {
	input := []byte{0, 0, 0, 7}
	input = append(input, []byte("ssh-rsa")...)
	input = append(input, []byte("comment and padding")...)
	a, ok := extractKeyAlgorithm(input)
	s.True(ok)
	s.Equal("ssh-rsa", a)

	input = []byte{0, 0, 0, 9}
	input = append(input, []byte("ssh-ecdsa")...)
	input = append(input, []byte("comment and padding")...)
	a, ok = extractKeyAlgorithm(input)
	s.True(ok)
	s.Equal("ssh-ecdsa", a)
}

func (s *sshSuite) Test_extractDummyCheckSum_FromAnEmptyInputItIsNotPossibleToExtractADummyCheckSum() {
	input := []byte{}

	_, _, ok := extractDummyCheckSum(input)

	s.False(ok)
}

func (s *sshSuite) Test_extractDummyCheckSum_ANotLongEnoughByteSliceShouldReturnNotOk() {
	input := []byte{1, 8, 10, 12, 16, 15}

	_, _, ok := extractDummyCheckSum(input)

	s.False(ok)
}

func (s *sshSuite) Test_extractDummyCheckSum_AnInputWith64BitsIsAValidDummyCheckSum() {
	input := []byte{1, 8, 10, 12, 16, 15, 11, 9}

	c, _, ok := extractDummyCheckSum(input)

	s.True(ok)
	s.Equal([]byte{1, 8, 10, 12, 16, 15, 11, 9}, c)
}

func (s *sshSuite) Test_extractDummyCheckSum_ALongEnoughInputReturnsAValidDummyCheckSum() {
	input := []byte{0, 3, 9, 13, 16, 18, 20, 22, 26, 2}

	c, r, ok := extractDummyCheckSum(input)

	s.True(ok)
	s.Equal([]byte{0, 3, 9, 13, 16, 18, 20, 22}, c)
	s.Equal([]byte{26, 2}, r)
}

func (s *sshSuite) Test_createPrivateKeyFrom_getsTheAlgorithmFromABasicKeyStructure() {
	input := []byte{
		0, 0, 0, 4,
		'n', 'o', 'n', 'e', // ciphername

		0, 0, 0, 0,
		// kdfname

		0, 0, 0, 0,
		// kdf

		0, 0, 0, 1, // number of keys

		0, 0, 0, 12, // ssh pub key name length

		0, 0, 0, 0,
		// public key keytype

		0, 0, 0, 0,
		// pub0

		0, 0, 0, 0,
		// pub1

		0, 0, 0, 17,

		0, 0, 0, 0, 0, 0, 0, 0, // checksum

		0, 0, 0, 5,
		byte('s'), byte('s'), byte('h'), byte('f'), byte('o'),
	}

	pk, ok := createPrivateKeyFrom(input)

	s.True(ok)
	s.Equal("sshfo", pk.algorithm)
}

func (s *sshSuite) Test_createPrivateKeyFrom_getsTheAlgorithmFromABasicKeyStructureWithOtherValues() {
	input := []byte{
		0, 0, 0, 4,
		'n', 'o', 'n', 'e', // ciphername

		0, 0, 0, 2,
		52, 42, // kdfname

		0, 0, 0, 1,
		7, // kdf

		0, 0, 0, 1, // number of keys

		0, 0, 0, 22, // ssh pub key name length

		0, 0, 0, 2,
		0, 0, // public key keytype

		0, 0, 0, 3,
		1, 2, 3, // pub0

		0, 0, 0, 5,
		7, 8, 12, 13, 55, // pub1

		0, 0, 0, 14,

		0, 1, 2, 3, 4, 5, 6, 7, // checksum

		0, 0, 0, 2,
		byte('T'), byte('A'),
	}

	pk, ok := createPrivateKeyFrom(input)

	s.True(ok)
	s.Equal("TA", pk.algorithm)
	s.False(pk.passwordProtected, "The key generated should not be marked as password protected")
}

func (s *sshSuite) Test_createPrivateKeyFrom_ReturnsAnErrorWhenReadingCiphernameFails() {
	input := []byte{
		77, 0, 0, 3, // incorrect length <------ here the error /!\
		1, 2, 3, // ciphername
	}

	_, ok := createPrivateKeyFrom(input)

	s.False(ok, "the ciphername length is not valid")
}

func (s *sshSuite) Test_createPrivateKeyFrom_ReturnsAnErrorWhenReadingKdfnameFails() {
	input := []byte{
		0, 0, 0, 3,
		1, 2, 3, // ciphername

		77, 0, 0, 2, // incorrect length <------ here the error /!\
		52, 42, // kdfname
	}

	_, ok := createPrivateKeyFrom(input)

	s.False(ok, "the kdfname length is not valid")
}

func (s *sshSuite) Test_createPrivateKeyFrom_ReturnsAnErrorWhenReadingKdfFails() {
	input := []byte{
		0, 0, 0, 3,
		1, 2, 3, // ciphername

		0, 0, 0, 2,
		52, 42, // kdfname

		77, 0, 0, 1, // incorrect length <------ here the error /!\
		7, // kdf
	}

	_, ok := createPrivateKeyFrom(input)

	s.False(ok, "the kdf length is not valid")
}

func (s *sshSuite) Test_createPrivateKeyFrom_ReturnsAnErrorWhenReadingTheNumberOfKeysFails() {
	input := []byte{
		0, 0, 0, 3,
		1, 2, 3, // ciphername

		0, 0, 0, 2,
		52, 42, // kdfname

		0, 0, 0, 1,
		7, // kdf

		0, 0, 0, // lacking 1 byte <------ here the error /!\
	}

	_, ok := createPrivateKeyFrom(input)

	s.False(ok, "32 bits were expected when reading the number of keys")
}

func (s *sshSuite) Test_createPrivateKeyFrom_IfTheNumberOfKeysIsZeroThenItIsNotAValidPrivateKey() {
	input := []byte{
		0, 0, 0, 3,
		1, 2, 3, // ciphername

		0, 0, 0, 2,
		52, 42, // kdfname

		0, 0, 0, 1,
		7, // kdf

		0, 0, 0, 0, // 0 keys  <------ here the error /!\
	}

	_, ok := createPrivateKeyFrom(input)

	s.False(ok, "at least one key is expected")
}

func (s *sshSuite) Test_createPrivateKeyFrom_ReturnsAnErrorWhenReadingPublicKeyBlockFails() {
	input := []byte{
		0, 0, 0, 3,
		1, 2, 3, // ciphername

		0, 0, 0, 2,
		52, 42, // kdfname

		0, 0, 0, 1,
		7, // kdf

		0, 0, 0, 1, // number of keys

		77, 0, 0, 22, // incorrect length <------ here the error /!\

		0, 0, 0, 2,
		0, 0, // public key keytype

	}

	_, ok := createPrivateKeyFrom(input)

	s.False(ok, "the public key block length is not valid")
}

func (s *sshSuite) Test_createPrivateKeyFrom_ReturnsAnErrorWhenReadingPrivateKeyBlockFails() {
	input := []byte{
		0, 0, 0, 3,
		1, 2, 3, // ciphername

		0, 0, 0, 2,
		52, 42, // kdfname

		0, 0, 0, 1,
		7, // kdf

		0, 0, 0, 1, // number of keys

		0, 0, 0, 22, // ssh pub key name length

		0, 0, 0, 2,
		0, 0, // public key keytype

		0, 0, 0, 3,
		1, 2, 3, // pub0

		0, 0, 0, 5,
		7, 8, 12, 13, 55, // pub1

		77, 0, 0, 14, // incorrect length <------ here the error /!\

		0, 1, 2, 3, 4, 5, 6, 7, // checksum

		0, 0, 0, 2,
		byte('T'), byte('A'),
	}

	_, ok := createPrivateKeyFrom(input)

	s.False(ok, "the private key block length is not valid")
}

func (s *sshSuite) Test_createPrivateKeyFrom_ReturnsAnErrorWhenReadingDummyChecksumFails() {
	input := []byte{
		0, 0, 0, 4,
		'n', 'o', 'n', 'e', // ciphername

		0, 0, 0, 2,
		52, 42, // kdfname

		0, 0, 0, 1,
		7, // kdf

		0, 0, 0, 1, // number of keys

		0, 0, 0, 22, // ssh pub key name length

		0, 0, 0, 2,
		0, 0, // public key keytype

		0, 0, 0, 3,
		1, 2, 3, // pub0

		0, 0, 0, 5,
		7, 8, 12, 13, 55, // pub1

		0, 0, 0, 6,

		0, 1, 2, 3, 4, 5, // incorrect dummy checksum <------ here the error /!\
	}

	_, ok := createPrivateKeyFrom(input)

	s.False(ok, "the dummy checksum length is not valid")
}

func (s *sshSuite) Test_createPrivateKeyFrom_ReturnsAnErrorWhenExtractingPrivateKeyAlgorithmFails() {
	input := []byte{
		0, 0, 0, 4,
		'n', 'o', 'n', 'e', // ciphername

		0, 0, 0, 2,
		52, 42, // kdfname

		0, 0, 0, 1,
		7, // kdf

		0, 0, 0, 1, // number of keys

		0, 0, 0, 22, // ssh pub key name length

		0, 0, 0, 2,
		0, 0, // public key keytype

		0, 0, 0, 3,
		1, 2, 3, // pub0

		0, 0, 0, 5,
		7, 8, 12, 13, 55, // pub1

		0, 0, 0, 13,

		0, 1, 2, 3, 4, 5, 6, 7, // checksum

		0, 0, 0, 2,
		byte('T'),
	}

	_, ok := createPrivateKeyFrom(input)

	s.False(ok, "the private key algorithm is not valid")
}

const correctRSAPasswordProtectedKey = `
-----BEGIN OPENSSH PRIVATE KEY-----
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
I7kWTQ==
-----END OPENSSH PRIVATE KEY-----`

func (s *sshSuite) Test_createPrivateKeyFrom_WillReturnOKForAValidPasswordProtectedKey() {
	// This key was taken from https://peterlyons.com/problog/2017/12/openssh-ed25519-private-key-file-format/
	// It is encrypted using aes256-cbc with bcrypt as the KDF
	validKey, _ := hex.DecodeString("6f70656e7373682d6b65792d7631000000000a6165733235362d63626300" +
		"0000066263727970740000001800000010d08f6b8fd17593f246db4ac6c4" +
		"5a11930000001000000001000000330000000b7373682d65643235353139" +
		"0000002062837be86c63712896b8e0e7543e367c3abd0c0b5ad3e764ea0e" +
		"4f8ddd7d00ef000000901e60c56ef30d0ff02e07b57bf14645076c32c86c" +
		"88ecad545ca28424e4739aff5895bebd6778e70b6c54b309b9fdb0c94102" +
		"bf8cef5b97d3d75636967e67e4b9c1ee72ae81074b0ce0f7e540e051d569" +
		"05da263af3e383342cc75b3145242abb75257586a119c9d3673dfb7eabe4" +
		"696350904e7c7af3cd77f28bea10374e15bc6536c2e1029438fdd3930bee" +
		"bbc5ac30")
	validKey = bytes.TrimPrefix(validKey, privateKeyAuthMagicWithTerminator)

	pk, ok := createPrivateKeyFrom(validKey)
	s.True(ok, "A password protected private key should be possible to parse")
	s.True(pk.passwordProtected, "The key generated should be marked as password protected")
	s.Equal("ssh-ed25519", pk.algorithm)
}

func (s *sshSuite) Test_isRSAPrivateKey_CheckIfAStringHasTheFormatOfAnRSAPrivateKey() {
	a, _ := accessWithTestLogging()
	k := ""
	s.False(a.isRSAPrivateKey(k), "An empty string is not an OpenSSH private key representation thus it is not an RSA key")

	k = correctRSASSHPrivateKey
	s.True(a.isRSAPrivateKey(k), "A string with the algorithm identifier ssh-rsa is an RSA key")

	k = correctECDSASSHPrivateKey
	s.False(a.isRSAPrivateKey(k), "A string with the algorithm identifier ssh-ecdsa is not an RSA key")
}
