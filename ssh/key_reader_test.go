package ssh

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
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
		file := filepath.Join(s.tdir, f)
		err := os.WriteFile(file, []byte("some content"), 0666)

		s.Nil(err)
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

	k = "ssh-rsa bla batman@debian"
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

	k = "ssh-rsa  AAAAA foo@debian"
	_, ok = parsePublicKey(k)
	s.True(ok, "More than one whitespace character serves as one single separator between columns")

	k = "ssh-rsa\tAAAAA foo@debian"
	_, ok = parsePublicKey(k)
	s.True(ok, "A tab can be a separator for columns")

	k = "ssh-rsa   \t  \t  \t AAAAA foo@debian"
	pub, ok = parsePublicKey(k)
	s.True(ok, "A mix of tabs and spaces serve as one single separator")
	s.Equal("AAAAA", pub.key)

	k = "ssh-rsa AAQQ foo@debian foo2@debian"
	pub, ok = parsePublicKey(k)
	s.True(ok, "More than one comment is acceptable in an SSH public key")
	s.Equal("foo@debian foo2@debian", pub.comment)

	k = "ssh-rsa AAQQ"
	_, ok = parsePublicKey(k)
	s.True(ok, "An SSH public key without a comment is still acceptable")
}

func (s *sshSuite) Test_CheckIfThePublicKeyTypeIdentifierIsRSA() {
	pub := publicKey{}
	s.False(pub.isRSA(), "An empty key is not an RSA key")

	pub = publicKey{algorithm: rsaAlgorithm}
	s.True(pub.isRSA(), "A key with the algorithm identifier ssh-rsa is an RSA key")

	pub = publicKey{algorithm: "ssh-ecdsa"}
	s.False(pub.isRSA(), "A key with the algorithm identifier ssh-ecdsa is not an RSA key")
}

func (s *sshSuite) Test_CheckIfAStringHasTheFormatOfAnRSAPublicKey() {
	k := ""
	s.False(isRSAPublicKey(k), "An empty string is not an SSH public key representation thus it is not an RSA key")

	k = "ssh-rsa AAQQ"
	s.True(isRSAPublicKey(k), "A string with the algorithm identifier ssh-rsa is an RSA key")

	k = "ssh-ecdsa AAQQ"
	s.False(isRSAPublicKey(k), "A string with the algorithm identifier ssh-ecdsa is not an RSA key")
}

func (s *sshSuite) Test_parsePrivateKey_AnEmptyStringIsNotAPrivateKey() {
	pk := ""

	_, ok := parsePrivateKey(pk)

	s.False(ok)
}

func (s *sshSuite) Test_parsePrivateKey_AStringThatDoesNotFollowThePEMFormatShouldNotBeConsideredAPrivateKey() {
	pk := "This is not a RSA Private Key and this does not follow the PEM format"

	_, ok := parsePrivateKey(pk)

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

	_, ok := parsePrivateKey(pk)

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

	_, ok := parsePrivateKey(pk)

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

	_, ok := parsePrivateKey(pk)

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

	_, ok := parsePrivateKey(pk)

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

	priv, ok := parsePrivateKey(pk)

	s.True(ok)
	s.Equal("ecdsa-sha2-nistp256", priv.algorithm)
}

func corruptMagicValue(s string) string {
	return strings.ReplaceAll(s, sshMagicValueEncoded, "b3BlbnLzaC1rZXktdjE")
}

func (s *sshSuite) Test_parsePrivateKey_AStringContainingABinaryWithoutACorrectSSHMagicValueIsNotAValidPrivateKey() {
	pk := corruptMagicValue(correctECDSASSHPrivateKey)

	_, ok := parsePrivateKey(pk)

	s.False(ok)
}

func (s *sshSuite) Test_parsePrivateKey_AStringContainingAVeryShortBase64StringIsNotAValidPrivateKey() {
	pk := `
-----BEGIN OPENSSH PRIVATE KEY-----
b3Bl
-----END OPENSSH PRIVATE KEY-----
`
	_, ok := parsePrivateKey(pk)

	s.False(ok)
}

func (s *sshSuite) Test_parsePrivateKey_AStringContainingAWellFormedRSAOpenSSHPrivateKeyShouldBeConsideredAPrivateKey() {
	pk := `
-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----
`

	priv, ok := parsePrivateKey(pk)

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
	_, ok := extractPrivateKeyAlgorithm(input)
	s.False(ok)
}

func (s *sshSuite) Test_extractPrivateKeyAlgorithm_ANotLongEnoughByteSliceShouldReturnNotOk() {
	input := []byte{0, 0, 0, 7}
	input = append(input, []byte("ssh-r")...)

	_, ok := extractPrivateKeyAlgorithm(input)
	s.False(ok)
}

func (s *sshSuite) Test_extractPrivateKeyAlgorithm_ALongEnoughByteSliceShouldReturnAValidPrivateKeyAlgorithm() {
	input := []byte{0, 0, 0, 7}
	input = append(input, []byte("ssh-rsa")...)
	input = append(input, []byte("comment and padding")...)
	a, ok := extractPrivateKeyAlgorithm(input)
	s.True(ok)
	s.Equal("ssh-rsa", a)

	input = []byte{0, 0, 0, 9}
	input = append(input, []byte("ssh-ecdsa")...)
	input = append(input, []byte("comment and padding")...)
	a, ok = extractPrivateKeyAlgorithm(input)
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
		0, 0, 0, 0,
		// ciphername

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

		0, 0, 0, 14,

		0, 1, 2, 3, 4, 5, 6, 7, // checksum

		0, 0, 0, 2,
		byte('T'), byte('A'),
	}

	pk, ok := createPrivateKeyFrom(input)

	s.True(ok)
	s.Equal("TA", pk.algorithm)
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

		0, 0, 0, 6,

		0, 1, 2, 3, 4, 5, // incorrect dummy checksum <------ here the error /!\
	}

	_, ok := createPrivateKeyFrom(input)

	s.False(ok, "the dummy checksum length is not valid")
}

func (s *sshSuite) Test_createPrivateKeyFrom_ReturnsAnErrorWhenExtractingPrivateKeyAlgorithmFails() {
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

		0, 0, 0, 13,

		0, 1, 2, 3, 4, 5, 6, 7, // checksum

		0, 0, 0, 2,
		byte('T'),
	}

	_, ok := createPrivateKeyFrom(input)

	s.False(ok, "the private key algorithm is not valid")
}
