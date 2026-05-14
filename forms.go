package forms

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jaredtmartin/bolt-go"
	"github.com/jaredtmartin/hound"
)

type RenderMethod string

const (
	RenderByName      RenderMethod = "RenderByName"
	RenderByFieldType RenderMethod = "RenderByFieldType"
	RenderByDataType  RenderMethod = "RenderByDataType"
)

func (m RenderMethod) String() string {
	return string(m)
}

type FieldTypes string

func (t FieldTypes) String() string {
	return string(t)
}

const (
	IdFieldType       FieldTypes = "Id"
	TextFieldType     FieldTypes = "Text"
	HiddenFieldType   FieldTypes = "Hidden"
	CheckboxFieldType FieldTypes = "Checkbox"
	NumberFieldType   FieldTypes = "Number"
	EmailFieldType    FieldTypes = "Email"
	PhoneFieldType    FieldTypes = "Phone"
	// TagsFieldType     FieldTypes = "Tags"
	RadioFieldType    FieldTypes = "Radio"
	TextareaFieldType FieldTypes = "Textarea"
	SelectFieldType   FieldTypes = "Select"
)

type Model struct {
	hound.Model     `bson:",inline"`
	fieldConfig     map[string]*FieldConfig `json:"-" bson:"-"`
	fieldComponents map[string]Component    `json:"-" bson:"-"`
	formatters      map[string]Formatter    `json:"-" bson:"-"`
}
type FieldConfig struct {
	name      string
	label     string
	component Component
	formatter Formatter
}

// type FieldBuilder struct {
// 	config *FieldConfig
// }

func NewModel(collectionName string, id ...string) Model {
	return Model{
		Model: hound.NewModel(collectionName, id...),
		fieldConfig: map[string]*FieldConfig{
			"Model": {component: HiddenIdField},
		},
		formatters: map[string]Formatter{
			"int":   IntFormatter,
			"int32": IntFormatter,
			"int64": IntFormatter,
			"bool":  BoolFormatter,
		},
		fieldComponents: map[string]Component{
			TextFieldType.String():     TextField,
			HiddenFieldType.String():   HiddenField,
			IdFieldType.String():       HiddenIdField,
			TextareaFieldType.String(): TextareaField,
			NumberFieldType.String():   NumberField,
			EmailFieldType.String():    EmailField,
			PhoneFieldType.String():    PhoneField,
			// TagsFieldType.String():     defaultTagsField,
			CheckboxFieldType.String(): CheckboxField,
			RadioFieldType.String():    RadioField,
			SelectFieldType.String():   SelectField,
		},
	}
}

// Three ways to override fields:
//  1. Have a custom render method on the struct for that field
//  2. Provide a tag to specify what type of field you want.
//  3. Have a custom render method on the value type

type Component func(name, label, value string) *bolt.Field
type Formatter func(value reflect.Value) string

func (m *Model) UseComponent(fieldType string, component Component) {
	m.fieldComponents[fieldType] = component
}
func (m *Model) FieldConfig(name string) *FieldConfig {
	config := &FieldConfig{}
	m.fieldConfig[name] = config
	return config
}
func (c *FieldConfig) Name(name string) {
	c.name = name
}
func (c *FieldConfig) Label(label string) {
	c.label = label
}
func (c *FieldConfig) Component(component Component) {
	c.component = component
}
func (c *FieldConfig) Formatter(formatter Formatter) {
	c.formatter = formatter
}

func (m *Model) getComponent(meta reflect.StructField) Component {
	// Check if there's a renderer specified for the field by name

	if config, ok := m.fieldConfig[meta.Name]; ok {
		if config.component != nil {
			return config.component
		}
	}

	// log.Println(`meta.Tag.Get("element"): `, meta.Tag.Get("element"))
	if tagEl := meta.Tag.Get("element"); tagEl != "" {
		if el, ok := m.fieldComponents[tagEl]; ok {
			return el
		}
	}
	if elName := getElementFromDataType(meta.Type.String()); elName != "" {
		if el, ok := m.fieldComponents[elName.String()]; ok {
			return el
		}
	}
	return m.fieldComponents[TextFieldType.String()]
}
func blankField() *bolt.Field {
	return &bolt.Field{DefaultElement: bolt.NewDefaultElement("")}
}
func getReflectTypeAndValue(obj any) (reflect.Type, reflect.Value) {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	// Dereference if pointer
	if objType.Kind() == reflect.Pointer {
		objType = objType.Elem()
		objValue = objValue.Elem()
	}
	return objType, objValue
}
func (m *Model) renderField(meta reflect.StructField, value reflect.Value, prefix ...string) *bolt.Field {
	namePrefix := ""
	if len(prefix) > 0 {
		namePrefix = prefix[0]
	}
	component := m.getComponent(meta)
	if component == nil {
		log.Println(`getComponent returned no component!`)
	}
	return component(namePrefix+m.getName(meta), m.getLabel(meta), m.getValue(meta, value))
}
func (m *Model) Field(name string, obj any) *bolt.Field {
	objType, objValue := getReflectTypeAndValue(obj)
	// Get the struct field by name (on the TYPE, not the value)
	meta, ok := objType.FieldByName(name)
	if !ok {
		log.Printf("Field with name %s not found\n", name)
		return blankField()
	}
	// Get the value
	value := objValue.FieldByName(name)
	if !value.IsValid() {
		return blankField()
	}
	return m.renderField(meta, value)
	// log.Printf("Field with name %s not found\n", name)
	// return &bolt.Field{DefaultElement: bolt.NewDefaultElement("")}
}
func (m *Model) getFormatter(meta reflect.StructField) Formatter {
	if config, ok := m.fieldConfig[meta.Name]; ok && config.formatter != nil {
		return config.formatter
	}
	format := meta.Tag.Get("format")
	if format != "" {
		if formatter, ok := m.formatters[format]; ok {
			return formatter
		}
	}
	return StringFormatter
}
func (m *Model) getValue(meta reflect.StructField, value reflect.Value) string {
	formatter := m.getFormatter(meta)
	log.Println(`formatter: `, formatter)
	return formatter(value)
}
func (m *Model) getLabel(meta reflect.StructField) string {
	config, ok := m.fieldConfig[meta.Name]
	if ok {
		return config.label
	}
	label := meta.Tag.Get("label")
	if label != "" {
		return label
	}
	return meta.Name
}

func (m *Model) getName(meta reflect.StructField) string {
	config, ok := m.fieldConfig[meta.Name]
	if ok {
		return config.name
	}
	name := meta.Tag.Get("name")
	if name != "" {
		return name
	}
	return meta.Name
}
func getElementFromDataType(datatype string) FieldTypes {
	switch datatype {
	case "string":
		return TextFieldType
	case "int":
		return NumberFieldType
	case "int32":
		return NumberFieldType
	case "int64":
		return NumberFieldType
	case "bool":
		return CheckboxFieldType
	default:
		return TextFieldType
	}
}
func (m *Model) Form(s any, prefix ...string) bolt.Element {
	objType, objValue := getReflectTypeAndValue(s)
	log.Printf("This object has %d fields to render. ", objType.NumField())
	form := bolt.Form()
	for i := 0; i < objType.NumField(); i++ {

		meta := objType.Field(i)
		value := objValue.Field(i)

		// Skip unexported fields (reflection can't access them)
		if !meta.IsExported() {
			log.Printf(`Skipping %s of type %s becuase it is not exported and reflection can't access it`, meta.Name, meta.Type)
			continue
		}
		log.Printf(`Rendering %s of type %s`, meta.Name, meta.Type)
		form.Add(m.renderField(meta, value, prefix...))
	}
	return form
}

func HiddenField(name, label, value string) *bolt.Field {
	log.Println("Rendering HiddenField with label", label)
	inputEl := bolt.HiddenInput(name, value)
	field := &bolt.Field{
		DefaultElement: bolt.NewDefaultElement(""),
		Input:          inputEl,
	}
	field.Children(inputEl)
	return field
}
func HiddenIdField(name, label, value string) *bolt.Field {
	log.Println("Rendering IdField with label", label)
	return HiddenField("Id", "", value)
}
func TextField(name, label, value string) *bolt.Field {
	log.Println("Rendering TextField with label", label)
	return bolt.TextField(name, label, value)
}
func TextareaField(name, label, value string) *bolt.Field {
	log.Println("Rendering TextareaField with label", label)
	return bolt.Textarea(name, label, value)
}
func NumberField(name, label, value string) *bolt.Field {
	log.Println("Rendering NumberField with label", label)
	input := bolt.TextField(name, label, value)
	input.Type("number")
	return input
}
func EmailField(name, label, value string) *bolt.Field {
	log.Println("Rendering EmailField with label", label)
	input := bolt.TextField(name, label, value)
	input.Type("email")
	return input
}
func PhoneField(name, label, value string) *bolt.Field {
	log.Println("Rendering PhoneField with label", label)
	input := bolt.TextField(name, label, value)
	input.Type("phone")
	return input
}

// func defaultTagsField(name, label, value string) *bolt.Field {

// }
func CheckboxField(name, label, value string) *bolt.Field {
	log.Println("Rendering CheckboxField with label", label)
	return bolt.Checkbox(name, label, value)
}
func RadioField(name, label, value string) *bolt.Field {
	log.Println("Rendering RadioField with label", label)
	return bolt.Radio(name, label, value)
}
func SelectField(name, label, value string) *bolt.Field {
	log.Println("Rendering SelectField with label", label)
	return bolt.Select(name, label, value, []bolt.Option{
		{Label: "Male", Value: "male"},
		{Label: "Female", Value: "female"},
	})
}
func StringFormatter(v reflect.Value) string {
	return fmt.Sprint(v.String())
}
func IntFormatter(v reflect.Value) string {
	return fmt.Sprint(v.Int())
}
func BoolFormatter(v reflect.Value) string {
	if v.Bool() {
		return "true"
	}
	return ""
}
func IdFormatter(v reflect.Value) string {
	if modelStruct, ok := v.Interface().(Model); ok {
		// fmt.Println("ID via assertion:", modelStruct.Id)
		return modelStruct.Id
	}
	log.Println("Unable to ")
	return ""
}
func StringSliceFormatter(v reflect.Value) string {
	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.String {
		return strings.Join(v.Interface().([]string), ",")
	}
	return fmt.Sprint(v.Interface())
}

func (m *Model) Element(name ...string) bolt.Element {
	nameFormat := "Id"
	if len(name) > 0 {
		nameFormat = fmt.Sprintf(name[0], nameFormat)
	}
	return bolt.HiddenInput(nameFormat, m.Id).Attr("special", "ed")
}
