package forms

import (
	"fmt"
	"log"
	"reflect"

	"github.com/jaredtmartin/bolt-go"
	"github.com/jaredtmartin/hound"
)

type FieldTypes string

const (
	TextFieldType     FieldTypes = "TextField"
	CheckboxFieldType FieldTypes = "CheckboxField"
	NumberFieldType   FieldTypes = "NumberField"
	EmailFieldType    FieldTypes = "EmailField"
	PhoneFieldType    FieldTypes = "PhoneField"
	// TagsFieldType     FieldTypes = "TagsField"
	RadioFieldType    FieldTypes = "RadioField"
	TextareaFieldType FieldTypes = "TextareaField"
	SelectFieldType   FieldTypes = "SelectField"
)

type Model struct {
	hound.Model    `bson:",inline"`
	fieldRenderers map[string]FieldRenderer
	// DefaultTextField     FieldRenderer
	// DefaultTextareaField FieldRenderer
	// DefaultNumberField   FieldRenderer
	// DefaultEmailField    FieldRenderer
	// DefaultPhoneField    FieldRenderer
	// // DefaultTagsField     FieldRenderer
	// DefaultCheckboxField FieldRenderer
	// DefaultRadioField    FieldRenderer
	// DefaultSelectField   FieldRenderer
}

func NewModel(collectionName string, id ...string) Model {
	return Model{
		Model: hound.NewModel(collectionName, id...),
		fieldRenderers: map[string]FieldRenderer{
			"TextField":     defaultTextField,
			"TextareaField": defaultTextareaField,
			"NumberField":   defaultNumberField,
			"EmailField":    defaultEmailField,
			"PhoneField":    defaultPhoneField,
			// "TagsField":     defaultTagsField,
			"CheckboxField": defaultCheckboxField,
			"RadioField":    defaultRadioField,
			"SelectField":   defaultSelectField,
		},
	}
}

// Three ways to override fields:
//  1. Have a custom render method on the struct for that field
//  2. Provide a tag to specify what type of field you want.
//  3. Have a custom render method on the value type

type FieldRenderer func(name, label, value string) *bolt.Field

func (m *Model) SetRenderer(name string, renderer FieldRenderer) {
	m.fieldRenderers[name] = renderer
}
func (m *Model) GetRenderer(name string) FieldRenderer {
	return m.fieldRenderers[name]
}
func (m *Model) RenderField(name string, obj any) *bolt.Field {
	object := reflect.ValueOf(obj)
	if object.Kind() == reflect.Pointer {
		object = object.Elem()
	}
	value := object.FieldByName(name)
	if value.IsValid() {
		renderer := m.fieldRenderers[name]
		if renderer != nil {
			return renderer(name, name, value.String())
		}
	}
	log.Printf("Field with name %s not found\n", name)
	return &bolt.Field{DefaultElement: bolt.NewDefaultElement("")}
}
func (m *Model) getRendererFromObject(object reflect.Value, field reflect.StructField) (FieldRenderer, bool) {
	customRenderMethodName := fmt.Sprintf("Render%sField", field.Name)
	customRenderMethod := object.MethodByName(customRenderMethodName)
	if customRenderMethod.IsValid() {
		if renderer, ok := customRenderMethod.Interface().(FieldRenderer); ok {
			// log.Printf("Rendering %s of type %s using %s", field.Name, field.Type, customRenderMethodName)
			return renderer, true
		}
	}
	return nil, false
}

// func (m *Model) getRendererFromTag(field reflect.StructField) (FieldRenderer, bool) {
// 	if tag := field.Tag.Get("bolt"); tag != "" {
// 		switch tag {
// 		case string(TextFieldType):
// 			return m.DefaultTextField, true
// 		case string(CheckboxFieldType):
// 			return m.DefaultCheckboxField, true
// 		case string(NumberFieldType):
// 			return m.DefaultNumberField, true
// 		case string(EmailFieldType):
// 			return m.DefaultEmailField, true
// 		case string(PhoneFieldType):
// 			return m.DefaultPhoneField, true
// 		// case string(TagsFieldType):
// 		// 	return m.DefaultTagsField, true
// 		case string(RadioFieldType):
// 			return m.DefaultRadioField, true
// 		case string(TextareaFieldType):
// 			return m.DefaultTextareaField, true
// 		case string(SelectFieldType):
// 			return m.DefaultSelectField, true
// 		}
// 	}
// 	return nil, false
// }

// func (m *Model) getRendererFromField(value reflect.Value) (FieldRenderer, bool) {
// 	if value.Kind() == reflect.Struct && value.CanAddr() {
// 		method := value.Addr().MethodByName("Render")
// 		// log.Println(`method.IsValid(): `, method.IsValid())
// 		if method.IsValid() {
// 			fieldRendererType := reflect.TypeFor[FieldRenderer]()
// 			if method.Type().ConvertibleTo(fieldRendererType) {
// 				renderer := method.Convert(fieldRendererType).Interface().(FieldRenderer)
// 				return renderer, true
// 			}
// 		}
// 	}
// 	return nil, false
// }
// func (m *Model) getIdField(field reflect.StructField) (FieldRenderer, bool) {
// 	if field.Name != "Model" {
// 		return nil, false
// 	}
// 	return defaultIdField, true
// }
// func (m *Model) getFieldRenderer(object reflect.Value, field reflect.StructField, value reflect.Value) FieldRenderer {
// 	renderer, ok := m.getIdField(field)
// 	if ok {
// 		return renderer
// 	}
// 	renderer, ok = m.getRendererFromObject(object, field)
// 	if ok {
// 		return renderer
// 	}
// 	renderer, ok = m.getRendererFromTag(field)
// 	if ok {
// 		return renderer
// 	}
// 	renderer, ok = m.getRendererFromField(value)
// 	if ok {
// 		return renderer
// 	}

//		return m.DefaultTextField
//	}
// func (m *Model) Form(s any) bolt.Element {
// 	structType := reflect.TypeOf(s)
// 	object := reflect.ValueOf(s)

// 	// Ensure we are working with a struct
// 	if structType.Kind() == reflect.Pointer {
// 		structType = structType.Elem()
// 		object = object.Elem()
// 	}
// 	log.Printf("This object has %d fields to render. ", structType.NumField())
// 	form := bolt.Form()
// 	for i := 0; i < structType.NumField(); i++ {

//			field := structType.Field(i)
//			value := object.Field(i)
//			log.Printf(`Rendering %s of type %s`, field.Name, field.Type)
//			// Skip unexported fields (reflection can't access them)
//			if !field.IsExported() {
//				continue
//			}
//			renderer := m.getFieldRenderer(object, field, value)
//			form.Add(renderer(field.Name, field.Name, value.String()))
//		}
//		return form
//	}
func defaultIdField(name, label, value string) *bolt.Field {
	// NewElement("input").Value(value).Type("hidden").Name(name)
	return bolt.NewField("Id", "", value, "hidden")
	// nullElement := bolt.None()
	// bil
	// return &bolt.Field{
	// 	Label: nullElement,
	// Input: bolt.HiddenInput(name, value),
	// 	Error: nullElement,
	// 	Check: nullElement,
	// }
}
func defaultTextField(name, label, value string) *bolt.Field {
	return bolt.TextField(name, label, value)
}
func defaultTextareaField(name, label, value string) *bolt.Field {
	return bolt.Textarea(name, label, value)
}
func defaultNumberField(name, label, value string) *bolt.Field {
	input := bolt.TextField(name, label, value)
	input.Type("number")
	return input
}
func defaultEmailField(name, label, value string) *bolt.Field {
	input := bolt.TextField(name, label, value)
	input.Type("email")
	return input
}
func defaultPhoneField(name, label, value string) *bolt.Field {
	input := bolt.TextField(name, label, value)
	input.Type("phone")
	return input
}

// func defaultTagsField(name, label, value string) *bolt.Field {

// }
func defaultCheckboxField(name, label, value string) *bolt.Field {
	return bolt.Checkbox(name, label, value)
}
func defaultRadioField(name, label, value string) *bolt.Field {
	return bolt.Radio(name, label, value)
}
func defaultSelectField(name, label, value string) *bolt.Field {
	return bolt.Select(name, label, value, []bolt.Option{
		{Label: "Male", Value: "male"},
		{Label: "Female", Value: "female"},
	})
}

func (m *Model) Element(name ...string) bolt.Element {
	nameFormat := "Id"
	if len(name) > 0 {
		nameFormat = fmt.Sprintf(name[0], nameFormat)
	}
	return bolt.HiddenInput(nameFormat, m.Id).Attr("special", "ed")
}

// type FieldRenderer func(name ...string) bolt.Element

//	func Element(field reflect.StructField, value reflect.Value) bolt.Element {
//		switch field.Type.Kind() {
//		case reflect.Bool:
//			return atoms.Checkbox(field.Name, field.Name, value.Bool(), field.Name, "")
//		case reflect.Int:
//			return atoms.TextField(field.Name, field.Name, value.String()).Attr("type", "number")
//		default:
//			return atoms.TextField(field.Name, field.Name, value.String())
//		}
//	}
// func RenderForm(s any) bolt.Element {
// 	structType := reflect.TypeOf(s)
// 	structValue := reflect.ValueOf(s)

// 	// Ensure we are working with a struct
// 	if structType.Kind() == reflect.Pointer {
// 		structType = structType.Elem()
// 		structValue = structValue.Elem()
// 	}

// 	form := bolt.Form()
// 	for i := 0; i < structType.NumField(); i++ {
// 		field := structType.Field(i)
// 		value := structValue.Field(i)

// 		// Skip unexported fields (reflection can't access them)
// 		if !field.IsExported() {
// 			continue
// 		}
// 		log.Println(`field.Name: `, field.Name)
// 		log.Println(`field.Type: `, field.Type)
// 		log.Printf(`value: %v\n`, value)
// 		customRenderMethodName := fmt.Sprintf("Render%sField", field.Name)
// 		customRenderMethod := structValue.MethodByName(customRenderMethodName)
// 		if customRenderMethod.IsValid() {
// 			// The method exists; you can now call it
// 			// log.Printf(`%s found\n`, customRenderMethodName)
// 			if renderer, ok := customRenderMethod.Interface().(FieldRenderer); ok {
// 				// Use the custom renderer
// 				log.Printf("Rendering %s of type %s using %s", field.Name, field.Type, customRenderMethodName)
// 				form.Add(renderer(field.Name, value.String()))
// 				// form.Add(renderer.Element(field, value))
// 				continue
// 			} else {
// 				log.Printf("%s doesn't implement FieldRenderer", customRenderMethodName)
// 			}
// 		}

// 		// Extract custom info from struct tags
// 		if tag := field.Tag.Get("bolt"); tag != "" {
// 			log.Printf(`found a bolt tag: %s\n`, tag)
// 			// (Simplified parsing logic for the example)
// 			if tag == "email" {
// 				form.Add(atoms.TextField(field.Name, field.Name, value.String())).Attr("type", "email")
// 				continue
// 			}
// 			// switch tag {
// 			// case "email":
// 			// 	form.Add(atoms.TextField())

// 			// }
// 		}
// 		log.Println(`field.Type.String(): `, field.Type.String())
// 		// if field.Name == "Model" {
// 		// 	log.Println("On Id field")
// 		// 	if value.Kind() == reflect.Struct && value.CanAddr() {
// 		// 		if modelPtr, ok := value.Addr().Interface().(*hound.Model); ok {
// 		// 			// success
// 		// 			log.Println("SUCCESS")
// 		// 			form.Add(bolt.HiddenInput("Id", modelPtr.Id))
// 		// 			continue
// 		// 		} else {
// 		// 			log.Println(`ok: `, ok)
// 		// 		}
// 		// 	} else {
// 		// 		log.Println(`value.Kind() == reflect.Struct: `, value.Kind() == reflect.Struct)
// 		// 		log.Println(`value.CanAddr(): `, value.CanAddr())
// 		// 	}
// 		// 	log.Println("#########PASSED ############")
// 		// }
// 		// if value.Kind() == reflect.Struct && value.CanAddr() {
// 		// 	if modelPtr, ok := value.Addr().Interface().(*hound.Model); ok {
// 		// 		form.Add(modelPtr.Element(field, value))
// 		// 	}
// 		// }
// 		// log.Println(`value.Kind() == reflect.Struct: `, value.Kind() == reflect.Struct)
// 		// log.Println(`value.CanAddr(): `, value.CanAddr())
// 		if value.Kind() == reflect.Struct && value.CanAddr() {
// 			method := value.Addr().MethodByName("Element")
// 			log.Println(`method.IsValid(): `, method.IsValid())
// 			if method.IsValid() {
// 				fieldRendererType := reflect.TypeFor[FieldRenderer]()
// 				if method.Type().ConvertibleTo(fieldRendererType) {
// 					renderer := method.Convert(fieldRendererType).Interface().(FieldRenderer)
// 					form.Add(renderer(field.Name, value.String()))
// 					continue
// 				}
// 			}
// 		}
// 		// if value.Kind() == reflect.Struct && value.CanAddr() {
// 		// 	renderFieldMethod := value.Addr().MethodByName("Element")
// 		// 	if renderFieldMethod.IsValid() {
// 		// 		// you found *hound.Model.Element
// 		// 		form.Add(Element(field, value))
// 		// 		continue
// 		// 	}
// 		// }
// 		// 	renderFieldMethod := value.MethodByName("Element")
// 		// 	log.Println(`renderFieldMethod: `, renderFieldMethod)
// 		// 	if renderFieldMethod.IsValid() {
// 		// 		log.Println(`renderFieldMethod.IsValid`)
// 		// 		if renderer, ok := renderFieldMethod.Interface().(FieldRenderer); ok {
// 		// 			// Use the custom renderer
// 		// 			form.Add(renderer(field, value))
// 		// 			// form.Add(renderer.Element(field, value))
// 		// 			continue
// 		// 		} else {
// 		// 			log.Printf("%s doesn't implement FieldRenderer", customRenderMethodName)
// 		// 		}
// 		// 	} else {
// 		// 		log.Println(`renderFieldMethod Is Not Valid`)
// 		// 	}

// 		form.Add(Element(field, value))

// 		// 	// fmt.Fprintf(&html, "<label>%s</label><input type='%s' name='%s' value='%v'><br>\n",
// 		// 	// 	label, inputType, field.Name, v.Field(i).Interface())
// 	}
// 	return form
// }
