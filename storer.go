package mellivora

import (
	"reflect"
)

const ()

type Indexes map[string]Index

type Instructions interface {
	Type() string
	Indexes() Indexes
}

type anonymousInstructions struct {
	reflectType reflect.Type
	indexes     Indexes
}

func newInstruction(model interface{}) Instructions {
	i, ok := model.(Instructions)
	if ok {
		return i
	}

	typ := reflect.TypeOf(model)
	for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	instructions := &anonymousInstructions{
		reflectType: typ,
		indexes:     map[string]Index{},
	}

	if instructions.reflectType.Name() == "" {
		panic("invalid model, must be a named type")
	}

	if instructions.reflectType.Kind() != reflect.Struct {
		panic("invalid model, must be a struct")
	}

	numberOfFields := instructions.reflectType.NumField()
	for i := 0; i < numberOfFields; i++ {
		indexName, unqique := "", false

	}
}
