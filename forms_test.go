package forms_test

import (
	"testing"

	"github.com/jaredtmartin/bolt-go"
	forms "github.com/jaredtmartin/musketforms"

	"github.com/stretchr/testify/assert"
)

type Model struct {
	forms.Model
}

func NewModel(col string, id ...string) Model {
	return Model{Model: forms.NewModel(col, id...)}
}
func (m *Model) DefaultTextField(name, label, value string, errors ...string) *bolt.Field {
	return bolt.NewField(name, label, value, "text")
}

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

type Status string

const (
	Available Status = "Available"
	Adopted   Status = "Adopted"
	Medical   Status = "Medical"
)

type Dog struct {
	// Should render HiddenIdField by tag
	Model
	// Should render default TextField
	Name string
	// Should render NumberField by data type
	Age int `name:"dob"`
	// Should render NumberField by data type
	Value64 int64 `label:"64bit"`
	// Should render NumberField by data type
	Value32 int32
	// Should not render by Default
	Tags []string
	// Should render SelectField by Tag
	Gender Gender `element:"Select"`
	// Should be able to override renderer
	Status Status
}
type Cat struct {
	Model
	Name  string
	Breed string
}

func NewDog(id string) *Dog {
	return &Dog{Model: NewModel("dogs", id)}
}

func TestSimpleForm(t *testing.T) {
	wix := &Cat{
		Model: NewModel("cats", "cat"),
		Name:  "Wix",
		Breed: "Mix",
	}
	result := wix.Form(wix).Render()
	expected := `<form><input name="Id" type="hidden" value="&lt;forms_test.Model Value&gt;"><div><label for="Name-field">Name</label><input id="Name-field" name="Name" type="text" value="Wix"><div id="Name-field-error"></div></div><div><label for="Breed-field">Breed</label><input id="Breed-field" name="Breed" type="text" value="Mix"><div id="Breed-field-error"></div></div></form>`
	assert.Equal(t, expected, result, "should match")
}
func TestSimpleField(t *testing.T) {
	wix := &Cat{
		Model: NewModel("cats", "cat"),
		Name:  "Wix",
		Breed: "Mix",
	}
	result := wix.Field("Name", wix).Render()
	expected := `<div><label for="Name-field">Name</label><input id="Name-field" name="Name" type="text" value="Wix"><div id="Name-field-error"></div></div>`
	assert.Equal(t, expected, result, "should match")
}
func TestBasicFields(t *testing.T) {
	spot := &Dog{
		Model:  NewModel("dogs", "spot"),
		Name:   "Spot",
		Age:    5,
		Tags:   []string{"Fun", "Silly"},
		Gender: Male,
		Status: Available,
	}
	result := spot.Field("Name", spot).Render()
	expected := `<div><label for="Name-field">Name</label><input id="Name-field" name="Name" type="text" value="Spot"><div id="Name-field-error"></div></div>`
	assert.Equal(t, expected, result, "should match")
}

func TestHiddenFieldFromTag(t *testing.T) {
	spot := &Dog{
		Model:  NewModel("dogs", "spot"),
		Name:   "Spot",
		Age:    5,
		Tags:   []string{"Fun", "Silly"},
		Gender: Male,
		Status: Available,
	}
	result := spot.Field("Name", spot).Render()
	expected := `<div><label for="Name-field">Name</label><input id="Name-field" name="Name" type="text" value="Spot"><div id="Name-field-error"></div></div>`
	assert.Equal(t, expected, result, "should match")
}
