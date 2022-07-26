package gui

import (
	"crypto/sha1"
	"fmt"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/api"
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

const publicKeyPathLabel = "publicKeyPath"
const privateKeyPathLabel = "privateKeyPath"
const publicKeyRowName = "keyDetailsPublicKeyRow"
const privateKeyRowName = "keyDetailsPrivateKeyRow"

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

func (kd *keyDetails) displayLocations(keyLocations []string, pathLabelName, rowName string) {
	if keyLocations != nil {
		label := kd.builder.get(pathLabelName).(gtki.Label)
		label.SetLabel(keyLocations[0])
		label.SetTooltipText(keyLocations[0])
	} else {
		row := kd.builder.get(rowName).(gtki.Box)
		row.Hide()
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

const fingerprintRow string = "keyFingerprintRow"

func (kd *keyDetails) displayFingerprint(rowName string) {
	if pk, ok := kd.key.(api.PublicKeyEntry); ok {
		f := formatFingerprint(pk.WithDigestContent(returningSlice20(sha1.Sum)))
		label := kd.builder.get("fingerprint").(gtki.Label)
		label.SetLabel(f)
		label.SetTooltipText(f)
	} else {
		row := kd.builder.get(rowName).(gtki.Box)
		row.Hide()
	}
}

func (kd *keyDetails) display() {
	kd.displayLocations(kd.key.PublicKeyLocations(), publicKeyPathLabel, publicKeyRowName)
	kd.displayLocations(kd.key.PrivateKeyLocations(), privateKeyPathLabel, privateKeyRowName)
	kd.displayFingerprint(fingerprintRow)
	kd.setClassForKeyDetails()
}

func (u *ui) populateKeyDetails(key api.KeyEntry, into gtki.Box) {
	clearAllChildrenOf[gtki.Widget](into)
	b, builder := buildObjectFrom[gtki.Box](u, "KeyDetails")

	kd := newKeyDetails(builder, key, b)
	kd.display()

	into.Add(b)
}
