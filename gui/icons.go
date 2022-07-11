package gui

import (
	"io/fs"
	"io/ioutil"
	"os"
)

func createTempFileWithContent(c []byte) string {
	f, _ := ioutil.TempFile("", "")
	f.Write(c)
	f.Close()
	return f.Name()
}

func withDefinitionInTempFile(def string, f func(string)) {
	content, _ := fs.ReadFile(getDefinitions(), def)
	file := createTempFileWithContent(content)
	defer os.Remove(file)

	f(file)
}

func (u *ui) loadResourceDefinitions() {
	withDefinitionInTempFile(resourceDefinitionPath("icons"), u.makeCustomIconsAvailable)
}

func (u *ui) makeCustomIconsAvailable(f string) {
	r, _ := u.gio.LoadResource(f)
	u.gio.RegisterResource(r)
	u.gtk.IconThemeGetDefault().AddResourcePath(keymirrorApplicationResourceID)
}
