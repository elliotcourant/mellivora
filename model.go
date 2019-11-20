package mellivora

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"math/rand"
	"reflect"
	"strings"
)

var (
	_ Model    = &modelInfo{}
	_ Field    = &modelField{}
	_ FieldSet = &fieldSet{}
)

type (
	Model interface {
		ModelId() uint32
		Name() string
		Fields() FieldSet
		PrimaryKey() FieldSet
		UniqueConstraints() UniqueConstraintSet
	}

	Field interface {
		FieldId() uint32
		Name() string
		IsPrimaryKey() bool
		Reflection() reflect.StructField
	}

	FieldSet interface {
		GetAll() []Field
		GetById(fieldId uint32) Field
		GetByName(fieldName string) Field
	}

	UniqueConstraint interface {
		UniqueConstraintId() uint32
		Name() string
		Fields() FieldSet
	}

	UniqueConstraintSet interface {
		GetAll() []UniqueConstraint
		GetById(uniqueConstraintId uint32) UniqueConstraint
		GetByName(uniqueConstraintName string) UniqueConstraint
	}
)

type modelInfo struct {
	modelId    uint32
	name       string
	fields     FieldSet
	primaryKey FieldSet
}

func (m *modelInfo) UniqueConstraints() UniqueConstraintSet {
	panic("implement me")
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

	fieldId      uint32
	name         string
	isPrimaryKey bool
	reflection   reflect.StructField
}

func (m *modelField) IsPrimaryKey() bool {
	return m.isPrimaryKey
}

func (m *modelField) Reflection() reflect.StructField {
	return m.reflection
}

func (m *modelField) FieldId() uint32 {
	return m.fieldId
}

func (m *modelField) Name() string {
	return m.name
}

type fieldSet struct {
	model  Model
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

func getModelInfo(model interface{}) Model {
	typ := getBaseTypeOf(model)

	mInfo := &modelInfo{
		modelId: rand.Uint32(),
		name:    typ.Name(),
	}

	fields := &fieldSet{
		model:  mInfo,
		fields: make([]Field, 0),
	}

	primaryKey := &fieldSet{
		model:  mInfo,
		fields: make([]Field, 0),
	}

	numFields := typ.NumField()
	for i := 0; i < numFields; i++ {
		reflection := typ.Field(i)

		field := &modelField{
			model:        mInfo,
			fieldId:      rand.Uint32(),
			name:         reflection.Name,
			isPrimaryKey: false,
			reflection:   reflection,
		}

		flags := getFlags(reflection.Tag.Get("m"))

		for key, value := range flags {
			switch key {
			case "pk":
				field.isPrimaryKey = true
			case "unique":
				fmt.Println(value)
			}
		}

		if field.isPrimaryKey {
			primaryKey.fields = append(primaryKey.fields, field)
		}

		fields.fields = append(fields.fields, field)
	}

	mInfo.fields = fields
	mInfo.primaryKey = primaryKey

	return mInfo
}

func getFlags(tag string) map[string]string {
	flags := strings.Split(tag, ",")
	items := map[string]string{}
	for _, flag := range flags {
		split := strings.SplitN(flag, ":", 2)
		switch len(split) {
		case 1:
			items[split[0]] = ""
		case 2:
			items[split[0]] = split[1]
		}
	}
	return items
}
