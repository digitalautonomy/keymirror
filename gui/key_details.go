package gui

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/digitalautonomy/keymirror/api"
	"github.com/digitalautonomy/keymirror/i18n"
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
	addClass(kd.box, fmt.Sprintf("algorithm-%s", strings.ToLower(kd.key.Algorithm().Name())))
	if kd.key.Algorithm().HasKeySize() {
		addClass(kd.box, fmt.Sprintf("key-size-%d", kd.key.Size()))
	}
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

const algorithmIdentifier = "algorithm"

func formatKeyAlgorithm(k api.KeyEntry) string {
	algo := k.Algorithm()
	if algo.HasKeySize() {
		// TODO: this formatting probably needs to be i18n in the future
		return fmt.Sprintf("%s (%d bits)", algo.Name(), k.Size())
	}
	return fmt.Sprintf("%s", algo.Name())
}

func (kd *keyDetails) displayAlgorithm() {
	label := kd.builder.get(algorithmIdentifier).(gtki.Label)
	label.SetLabel(formatKeyAlgorithm(kd.key))
}

const userIDLabelIdentifier = "userIDLabel"
const userIDIdentifier = "userID"

func (kd *keyDetails) displayUserID() {
	if pk, ok := kd.key.(api.PublicKeyEntry); ok {
		label := kd.builder.get(userIDIdentifier).(gtki.Label)
		label.SetLabel(pk.UserID())
	} else {
		kd.hideAll(userIDLabelIdentifier, userIDIdentifier)
	}
}

const notificationIdentifier = "notification"

func (kd *keyDetails) displayNotification(message string) {
	label := kd.builder.get(notificationIdentifier).(gtki.Label)
	label.SetLabel(message)
}

func (kd *keyDetails) hideNotificationRow() {
	kd.hide(notificationIdentifier)
}

func (kd *keyDetails) displayPotentialNotification() {
	switch kd.key.KeyType() {
	case api.PublicKeyType:
		kd.displayNotification(i18n.Local("(no private key available)"))
	case api.PrivateKeyType:
		kd.displayNotification(i18n.Local("(no public key available)"))
	case api.PairKeyType:
		fallthrough
	default:
		kd.hideNotificationRow()
	}
}

func (kd *keyDetails) display() {
	kd.displayPotentialNotification()
	kd.displayLocations(kd.key.PublicKeyLocations(), publicKeyPath, publicKeyPathLabel)
	kd.displayLocations(kd.key.PrivateKeyLocations(), privateKeyPath, privateKeyPathLabel)
	kd.displayIsPasswordProtected()
	kd.displayAlgorithm()
	kd.displayUserID()
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
