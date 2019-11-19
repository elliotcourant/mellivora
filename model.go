package mellivora

import (
	"fmt"
	"reflect"
)

type ModelInfo struct {
	ModelId uint32
	Name    string
}

type ModelField struct {
	model *ModelInfo

	FieldId uint32
	Name    string
	Reflect reflect.StructField
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

// func getModelInfo(model interface{}) ModelInfo {
// 	typ := getBaseTypeOf(model)
// }
