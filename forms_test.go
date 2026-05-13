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
	Model  `bson:",inline"`
	Name   string
	Age    int
	Tags   []string
	Gender Gender
	Status Status
}

func NewDog(id string) *Dog {
	return &Dog{Model: NewModel("dogs", id)}
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
	result := spot.RenderField("Name", spot).Render()
	assert.Equalf(t, `<form><div><input id="Id-field" name="Id" type="hidden" value="&lt;forms_test.Model Value&gt;"><div id="Id-field-error"></div></div><div><label for="Name-field">Name</label><input id="Name-field" name="Name" type="text" value="Spot"><div id="Name-field-error"></div></div><div><label for="Age-field">Age</label><input id="Age-field" name="Age" type="text" value="&lt;int Value&gt;"><div id="Age-field-error"></div></div><div><label for="Tags-field">Tags</label><input id="Tags-field" name="Tags" type="text" value="&lt;[]string Value&gt;"><div id="Tags-field-error"></div></div><div><label for="Gender-field">Gender</label><input id="Gender-field" name="Gender" type="text" value="Male"><div id="Gender-field-error"></div></div><div><label for="Status-field">Status</label><input id="Status-field" name="Status" type="text" value="Available"><div id="Status-field-error"></div></div></form>`, result, "should match")

}
