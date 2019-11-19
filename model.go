package mellivora

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"reflect"
)

var (
	_ Model    = &modelInfo{}
	_ Field    = &modelField{}
	_ FieldSet = &fieldSet{}
)

type Model interface {
	ModelId() uint32
	Name() string
	Fields() FieldSet
	PrimaryKey() FieldSet
}

type Field interface {
	FieldId() uint32
	Name() string
}

type FieldSet interface {
	GetAll() []Field
	GetById(fieldId uint32) Field
	GetByName(fieldName string) Field
}

type modelInfo struct {
	modelId    uint32
	name       string
	fields     FieldSet
	primaryKey FieldSet
}

func (m *modelInfo) ModelId() uint32 {
	return m.modelId
}

func (m *modelInfo) Name() string {
	return m.name
}

func (m *modelInfo) Fields() FieldSet {
	return m.fields
}

func (m *modelInfo) PrimaryKey() FieldSet {
	return m.primaryKey
}

type modelField struct {
	model *modelInfo

	fieldId    uint32
	name       string
	reflection reflect.StructField
}

func (m *modelField) FieldId() uint32 {
	return m.fieldId
}

func (m *modelField) Name() string {
	return m.name
}

type fieldSet struct {
	model *modelInfo

	fields []Field
}

func (f *fieldSet) GetById(fieldId uint32) Field {
	return linq.From(f.fields).FirstWith(func(i interface{}) bool {
		field, ok := i.(Field)
		return ok && field.FieldId() == fieldId
	}).(Field)
}

func (f *fieldSet) GetByName(fieldName string) Field {
	return linq.From(f.fields).FirstWith(func(i interface{}) bool {
		field, ok := i.(Field)
		return ok && field.Name() == fieldName
	}).(Field)
}

func (f *fieldSet) GetAll() []Field {
	return f.fields
}

func getBaseTypeOf(model interface{}) reflect.Type {
	typ := reflect.TypeOf(model)

	// We need to make sure we are working with the base type.
	for {
		switch typ.Kind() {
		case reflect.Ptr, reflect.Array, reflect.Slice:
			typ = typ.Elem()
		case reflect.Struct:
			return typ
		default:
			panic(fmt.Sprintf("%T/%s is not a supported type/kind", model, typ.Kind()))
		}
	}
}

// func getModelInfo(model interface{}) modelInfo {
// 	typ := getBaseTypeOf(model)
// }
