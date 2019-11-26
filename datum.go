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
		Reflection() reflect.Value
		Keys() map[string][]byte
		Verify() map[string]bool
	}

	datumBuilderBase struct {
		model  Model
		value  reflect.Value
		datums map[string][]byte
		verify map[string]bool
	}
)

func newDatumBuilder(model Model, value reflect.Value) datumBuilder {
	return &datumBuilderBase{
		model:  model,
		value:  value,
		datums: map[string][]byte{},
		verify: map[string]bool{},
	}
}

func (d *datumBuilderBase) Model() Model {
	return d.model
}

func (d *datumBuilderBase) Value() interface{} {
	return d.value.Interface()
}

func (d *datumBuilderBase) Reflection() reflect.Value {
	return d.value
}

func (d *datumBuilderBase) Keys() map[string][]byte {
	if len(d.datums) > 0 {
		return d.datums
	}

	// Handle initial datum record.
	{
		primaryKeyValueBuf := buffers.NewBytesBuffer()
		for _, fieldInfo := range d.model.PrimaryKey().GetAll() {
			primaryKeyValueBuf.AppendReflection(d.value.FieldByIndex(fieldInfo.Reflection().Index))
		}

		datumKeyBuf := buffers.NewBytesBuffer()
		datumKeyBuf.AppendByte(datumKeyPrefix)
		datumKeyBuf.AppendUint32(d.model.ModelId())
		datumKeyBuf.AppendRaw(primaryKeyValueBuf.Bytes())

		datumValueBuf := buffers.NewBytesBuffer()
		for _, fieldInfo := range d.model.Fields().GetAll() {
			if fieldInfo.IsPrimaryKey() {
				continue
			}
			datumValueBuf.AppendReflection(d.value.FieldByIndex(fieldInfo.Reflection().Index))
		}

		d.datums[string(datumKeyBuf.Bytes())] = datumValueBuf.Bytes()
	}

	for _, uniqueConstraint := range d.model.UniqueConstraints().GetAll() {
		uniqueConstraintBuf := buffers.NewBytesBuffer()
		uniqueConstraintBuf.AppendByte(uniqueKeyPrefix)
		uniqueConstraintBuf.AppendUint32(d.model.ModelId())
		uniqueConstraintBuf.AppendUint32(uniqueConstraint.UniqueConstraintId())
		for _, fieldInfo := range uniqueConstraint.Fields().GetAll() {
			uniqueConstraintBuf.AppendReflection(d.value.FieldByIndex(fieldInfo.Reflection().Index))
		}
		uniqueConstraintKey := uniqueConstraintBuf.Bytes()

		d.datums[string(uniqueConstraintKey)] = make([]byte, 0)

		// Make sure that the unique key does not exist.
		d.verify[string(uniqueConstraintKey)] = false
	}

	return d.datums
}

func (d *datumBuilderBase) Verify() map[string]bool {
	// If we have already built our datum set then we know the verify set has been built.
	if len(d.datums) > 0 {
		return d.verify
	}

	// If the datum set has not been built then the verify set is definitely not built.
	// Build the datum set and then return the verify map.
	_ = d.Keys()

	return d.verify
}
