package mellivora

import (
	"github.com/elliotcourant/buffers"
	"reflect"
)

var (
	_ datumBuilder = &datumBuilderBase{}
)

type (
	datumBuilder interface {
		Model() Model
		Value() interface{}
		Reflection() *reflect.Value
		Keys() map[string][]byte
		Verify() map[string][]byte
	}

	datumBuilderBase struct {
		model        Model
		value        *reflect.Value
		datums       map[string][]byte
		mustNotExist map[string]interface{}
	}
)

func newDatumBuilder(model Model, value *reflect.Value) datumBuilder {
	return &datumBuilderBase{
		model:        model,
		value:        value,
		datums:       map[string][]byte{},
		mustNotExist: map[string]interface{}{},
	}
}

func (d *datumBuilderBase) Model() Model {
	return d.model
}

func (d *datumBuilderBase) Value() interface{} {
	return d.value.Interface()
}

func (d *datumBuilderBase) Reflection() *reflect.Value {
	return d.value
}

func (d *datumBuilderBase) Keys() map[string][]byte {
	if len(d.datums) > 0 {
		return d.datums
	}

	primaryKeyValueBuf := buffers.NewBytesBuffer()
	primaryKeyBuf := buffers.NewBytesBuffer()
	primaryKeyBuf.AppendByte(datumKeyPrefix)
	primaryKeyBuf.AppendUint32(d.model.ModelId())
	for _, fieldInfo := range d.model.PrimaryKey().GetAll() {
		primaryKeyValueBuf.AppendReflection(d.value.FieldByIndex(fieldInfo.Reflection().Index))
	}
	primaryKeyBuf.AppendRaw(primaryKeyValueBuf.Bytes())
	d.datums[string(primaryKeyBuf.Bytes())] = make([]byte, 0)
	d.mustNotExist[string(primaryKeyBuf.Bytes())] = nil

	for _, fieldInfo := range d.model.Fields().GetAll() {
		if fieldInfo.IsPrimaryKey() {
			continue
		}
		fieldBuf := buffers.NewBytesBuffer()
		fieldBuf.AppendByte(datumKeyPrefix)
		fieldBuf.AppendUint32(d.model.ModelId())
		fieldBuf.AppendRaw(primaryKeyValueBuf.Bytes())
		fieldBuf.AppendUint32(fieldInfo.FieldId())

		valueBuf := buffers.NewBytesBuffer()
		valueBuf.AppendReflection(d.value.FieldByIndex(fieldInfo.Reflection().Index))
		d.datums[string(fieldBuf.Bytes())] = valueBuf.Bytes()
	}

	return d.datums
}

func (d *datumBuilderBase) Verify() map[string][]byte {
	panic("implement me")
}
