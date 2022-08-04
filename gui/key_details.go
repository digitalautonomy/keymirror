package gui

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/api"
	"strconv"
	"strings"
)

type clearable[T any] interface {
	GetChildren() []T
	Remove(T)
}

func clearAllChildrenOf[T any](b clearable[T]) {
	for _, c := range b.GetChildren() {
		b.Remove(c)
	}
}

type keyDetails struct {
	builder *builder
	key     api.KeyEntry
	box     gtki.Box
}

func newKeyDetails(builder *builder, key api.KeyEntry, box gtki.Box) *keyDetails {
	return &keyDetails{builder, key, box}
}

const publicKeyPath = "publicKeyPath"
const privateKeyPath = "privateKeyPath"
const publicKeyPathLabel = "publicKeyPathLabel"
const privateKeyPathLabel = "privateKeyPathLabel"

type hasStyleContext interface {
	GetStyleContext() (gtki.StyleContext, error)
}

func addClass(w hasStyleContext, class string) {
	sc, _ := w.GetStyleContext()
	sc.AddClass(class)
}

func removeClass(w hasStyleContext, class string) {
	sc, _ := w.GetStyleContext()
	sc.RemoveClass(class)
}

var keyTypeClassNames = map[api.KeyType]string{
	api.PublicKeyType:  "publicKey",
	api.PrivateKeyType: "privateKey",
	api.PairKeyType:    "keyPair",
}

func (kd *keyDetails) setClassForKeyDetails() {
	className := keyTypeClassNames[kd.key.KeyType()]
	addClass(kd.box, className)
}

type hideable interface {
	Hide()
}

func (kd *keyDetails) hideAll(ids ...string) {
	for _, id := range ids {
		kd.hide(id)
	}
}
func (kd *keyDetails) hide(id string) {
	l := kd.builder.get(id).(hideable)
	l.Hide()
}

func (kd *keyDetails) displayLocations(keyLocations []string, path, pathLabel string) {
	if keyLocations != nil {
		label := kd.builder.get(path).(gtki.Label)
		label.SetLabel(keyLocations[0])
		label.SetTooltipText(keyLocations[0])
	} else {
		kd.hideAll(pathLabel, path)
	}
}

func returningSlice20(f func([]byte) [20]byte) func([]byte) []byte {
	return func(v []byte) []byte {
		res := f(v)
		return res[:]
	}
}

func returningSlice32(f func([]byte) [32]byte) func([]byte) []byte {
	return func(v []byte) []byte {
		res := f(v)
		return res[:]
	}
}

func formatByteForFingerprint(b byte) string {
	return fmt.Sprintf("%02X", b)
}

func formatFingerprint(f []byte) string {
	result := []string{}

	for _, v := range f {
		result = append(result, formatByteForFingerprint(v))
	}

	return strings.Join(result, ":")
}

const sha1FingerprintLabel = "sha1FingerprintLabel"
const sha1Fingerprint = "sha1Fingerprint"
const sha256FingerprintLabel = "sha256FingerprintLabel"
const sha256Fingerprint = "sha256Fingerprint"

func (kd *keyDetails) displayFingerprint(fingerprintLabel, fingerprint string, f func([]byte) []byte) {
	if pk, ok := kd.key.(api.PublicKeyEntry); ok {
		fp := formatFingerprint(pk.WithDigestContent(f))
		label := kd.builder.get(fingerprint).(gtki.Label)
		label.SetLabel(fp)
		label.SetTooltipText(fp)
	} else {
		kd.hideAll(fingerprintLabel, fingerprint)
	}
}

const passwordProtectedLabel = "passwordProtectedLabel"

func (kd *keyDetails) privateKeyIsPasswordProtected() bool {
	pk, ok := kd.key.(api.PrivateKeyEntry)
	return ok && pk.IsPasswordProtected()
}

func (kd *keyDetails) displayIsPasswordProtected() {
	if kd.privateKeyIsPasswordProtected() {
		addClass(kd.box, "passwordProtectedPrivateKey")
	} else {
		kd.hide(passwordProtectedLabel)
	}
}

const sizeLabel = "sizeLabel"
const sizeValue = "size"
const sizeUnit = " bits"

func (kd *keyDetails) displaySize() {
	if pk, ok := kd.key.(api.KeyEntry); ok && pk.Size() > 0 {
		label := kd.builder.get(sizeValue).(gtki.Label)
		label.SetLabel(strconv.Itoa(pk.Size()) + sizeUnit)
	} else {
		kd.hideAll(sizeLabel, sizeValue)
	}
}

func (kd *keyDetails) display() {
	kd.displayLocations(kd.key.PublicKeyLocations(), publicKeyPath, publicKeyPathLabel)
	kd.displayLocations(kd.key.PrivateKeyLocations(), privateKeyPath, privateKeyPathLabel)
	kd.displayIsPasswordProtected()
	kd.displaySize()
	kd.displayFingerprint(sha1FingerprintLabel, sha1Fingerprint, returningSlice20(sha1.Sum))
	kd.displayFingerprint(sha256FingerprintLabel, sha256Fingerprint, returningSlice32(sha256.Sum256))
	kd.setClassForKeyDetails()
}

func (u *ui) populateKeyDetails(key api.KeyEntry, into gtki.Box) {
	clearAllChildrenOf[gtki.Widget](into)
	b, builder := buildObjectFrom[gtki.Box](u, "KeyDetails")

	kd := newKeyDetails(builder, key, b)
	kd.display()

	into.Add(b)
}
