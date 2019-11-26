package mellivora

import (
	"fmt"
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
		Keys() (map[string][]byte, error)
		Verify() (map[string]bool, error)
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

func (d *datumBuilderBase) setDatum(key, value []byte) error {
	if _, ok := d.datums[string(key)]; ok {
		return fmt.Errorf("an item with the key [%s] already exists in this datumset", string(key))
	}

	d.datums[string(key)] = value

	return nil
}

func (d *datumBuilderBase) setVerify(key []byte, canExist bool) error {
	if _, ok := d.verify[string(key)]; ok {
		return fmt.Errorf("an verify with the key [%s] already exists in this datumset", string(key))
	}

	d.verify[string(key)] = canExist

	return nil
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

func (d *datumBuilderBase) Keys() (map[string][]byte, error) {
	if len(d.datums) > 0 {
		return d.datums, nil
	}

	value := d.value

	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		numItems := value.Len()
		for i := 0; i < numItems; i++ {
			if err := d.encodeSingleDatum(value.Index(i)); err != nil {
				return nil, err
			}
		}
	default:
		if err := d.encodeSingleDatum(value); err != nil {
			return nil, err
		}
	}

	return d.datums, nil
}

func (d *datumBuilderBase) encodeSingleDatum(value reflect.Value) error {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Handle initial datum record.
	{
		primaryKeyValueBuf := buffers.NewBytesBuffer()
		for _, fieldInfo := range d.model.PrimaryKey().GetAll() {
			primaryKeyValueBuf.AppendReflection(value.FieldByIndex(fieldInfo.Reflection().Index))
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
			datumValueBuf.AppendReflection(value.FieldByIndex(fieldInfo.Reflection().Index))
		}

		if err := d.setDatum(datumKeyBuf.Bytes(), datumValueBuf.Bytes()); err != nil {
			return err
		}
	}

	for _, uniqueConstraint := range d.model.UniqueConstraints().GetAll() {
		uniqueConstraintBuf := buffers.NewBytesBuffer()
		uniqueConstraintBuf.AppendByte(uniqueKeyPrefix)
		uniqueConstraintBuf.AppendUint32(d.model.ModelId())
		uniqueConstraintBuf.AppendUint32(uniqueConstraint.UniqueConstraintId())
		for _, fieldInfo := range uniqueConstraint.Fields().GetAll() {
			uniqueConstraintBuf.AppendReflection(value.FieldByIndex(fieldInfo.Reflection().Index))
		}
		uniqueConstraintKey := uniqueConstraintBuf.Bytes()

		if err := d.setDatum(uniqueConstraintKey, make([]byte, 0)); err != nil {
			return err
		}

		// Make sure that the unique key does not exist.
		if err := d.setVerify(uniqueConstraintKey, false); err != nil {
			return err
		}
	}

	return nil
}

func (d *datumBuilderBase) Verify() (map[string]bool, error) {
	// If we have already built our datum set then we know the verify set has been built.
	if len(d.datums) > 0 {
		return d.verify, nil
	}

	// If the datum set has not been built then the verify set is definitely not built.
	// Build the datum set and then return the verify map.
	_, err := d.Keys()

	return d.verify, err
}
