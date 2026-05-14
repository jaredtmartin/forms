package forms_test

import (
	"strings"
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
	Age int `name:"dob" format:"int"`
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

func NewDog(name string) *Dog {
	return &Dog{
		Model: NewModel("dogs", name),
		Name:  strings.ToLower(name),
		Age:   5,
	}
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
func TestFormWithPrefix(t *testing.T) {
	wix := &Cat{
		Model: NewModel("cats", "cat"),
		Name:  "Wix",
		Breed: "Mix",
	}
	result := wix.Form(wix, "cats[0].").Render()
	expected := `<form><input name="Id" type="hidden" value="&lt;forms_test.Model Value&gt;"><div><label for="cats[0].Name-field">Name</label><input id="cats[0].Name-field" name="cats[0].Name" type="text" value="Wix"><div id="cats[0].Name-field-error"></div></div><div><label for="cats[0].Breed-field">Breed</label><input id="cats[0].Breed-field" name="cats[0].Breed" type="text" value="Mix"><div id="cats[0].Breed-field-error"></div></div></form>`
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
func TestFieldNameTag(t *testing.T) {
	spot := NewDog("Spot")
	result := spot.Field("Age", spot).Render()
	expected := `<div type="number"><label for="dob-field">Age</label><input id="dob-field" name="dob" type="text" value="5"><div id="dob-field-error"></div></div>`
	assert.Equal(t, expected, result, "should match")
}
func TestFieldNameOverride(t *testing.T) {
	spot := NewDog("Spot")
	spot.FieldConfig("Name").Name("Nickname")
	result := spot.Field("Name", spot).Render()
	expected := `<div><input id="Nickname-field" name="Nickname" type="text" value="spot"><div id="Nickname-field-error"></div></div>`
	assert.Equal(t, expected, result, "should match")
}
func TestIntValueFormatFromTag(t *testing.T) {
	spot := NewDog("Spot")
	result := spot.Field("Age", spot).Render()
	expected := `<div type="number"><label for="dob-field">Age</label><input id="dob-field" name="dob" type="text" value="5"><div id="dob-field-error"></div></div>`
	assert.Equal(t, expected, result, "should match")
}
func TestIntValueFormatFromConfig(t *testing.T) {
	spot := NewDog("Spot")
	spot.Value32 = 77
	result := spot.Field("Value32", spot).Render()
	expected := `<div type="number"><label for="Value32-field">Value32</label><input id="Value32-field" name="Value32" type="text" value="&lt;int32 Value&gt;"><div id="Value32-field-error"></div></div>`
	assert.Equal(t, expected, result, "should match")
	spot.FieldConfig("Value32").Formatter(forms.IntFormatter)
	result = spot.Field("Value32", spot).Render()
	expected = `<div type="number"><input id="-field" name type="text" value="77"><div id="-field-error"></div></div>`
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
