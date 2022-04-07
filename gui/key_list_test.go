package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3mocks/gtk"
	"github.com/digitalautonomy/keymirror/api"
	"github.com/stretchr/testify/mock"
)

type keyEntryMock struct {
	mock.Mock
}

func (ke *keyEntryMock) Locations() []string {
	returns := ke.Called()
	return returns.Get(0).([]string)
}

func (s *guiSuite) Test_createKeyEntryBoxFrom_CreatesAGTKIBoxWithTheGivenASSHKeyEntry() {
	box := s.setupBuildingOfKeyEntry("/home/amnesia/id_rsa.pub")

	u := &ui{gtk: s.gtkMock}

	keyEntry := &keyEntryMock{}
	keyEntry.On("Locations").Return([]string{"/home/amnesia/id_rsa.pub"}).Once()

	actualGtkBox := u.createKeyEntryBoxFrom(keyEntry)

	s.Equal(box, actualGtkBox)

	keyEntry.AssertExpectations(s.T())
}

type keyAccessMock struct {
	mock.Mock
}

func (ka *keyAccessMock) AllKeys() []api.KeyEntry {
	return ka.Called().Get(0).([]api.KeyEntry)
}

func fixedKeyAccess(keys ...api.KeyEntry) api.KeyAccess {
	ka := &keyAccessMock{}
	ka.On("AllKeys").Return(keys).Maybe()
	return ka
}

func fixedKeyEntry(location string) api.KeyEntry {
	ke := &keyEntryMock{}
	ke.On("Locations").Return([]string{location}).Maybe()
	return ke
}

func (s *guiSuite) setupBuildingOfObject(val interface{}, name string) {
	builder := &gtk.MockBuilder{}
	s.gtkMock.On("BuilderNew").Return(builder, nil).Once()
	builder.On("AddFromString", mock.Anything).Return(nil).Once()
	builder.On("GetObject", name).Return(val, nil).Once()
	s.addObjectToAssert(builder)
}

func (s *guiSuite) setupBuildingOfKeyEntry(path string) *gtk.MockBox {
	label := &gtk.MockLabel{}
	label.On("SetLabel", path).Return().Once()
	box := &gtk.MockBox{}
	box.On("GetChildren").Return([]gtki.Widget{label}).Once()
	s.setupBuildingOfObject(box, "KeyListEntry")
	s.addObjectToAssert(box)
	return box
}

func (s *guiSuite) Test_populateListWithKeyEntries_AddsAListOfKeyEntriesIntoAGTKBox() {
	ka := fixedKeyAccess(
		fixedKeyEntry("/home/amnesia/.ssh/id_rsa"),
		fixedKeyEntry("/home/amnesia/.ssh/id_ed25519"),
		fixedKeyEntry("/home/amnesia/.ssh/id_dsa"),
	)

	box1 := s.setupBuildingOfKeyEntry("/home/amnesia/.ssh/id_rsa")
	box2 := s.setupBuildingOfKeyEntry("/home/amnesia/.ssh/id_ed25519")
	box3 := s.setupBuildingOfKeyEntry("/home/amnesia/.ssh/id_dsa")

	box := &gtk.MockBox{}
	box.On("Add", box1).Return().Once()
	box.On("Add", box2).Return().Once()
	box.On("Add", box3).Return().Once()

	u := &ui{gtk: s.gtkMock}

	u.populateListWithKeyEntries(ka, box)

	box.AssertExpectations(s.T())
}
