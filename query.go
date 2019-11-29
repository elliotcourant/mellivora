package mellivora

import (
	"reflect"
	"strings"
)

type criteriaExpression = func(datum reflect.Value) bool

type Ex map[string]interface{}

type Query struct {
	destination reflect.Value
	model       Model
	txn         *Transaction
	filters     []Ex

	limit  int
	offset int
}

func (q *Query) Where(expression ...Ex) *Query {
	q.filters = append(q.filters, expression...)
	return q
}

func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

func (q *Query) Select(destination interface{}) error {
	dest := reflect.ValueOf(destination)
	for dest.Kind() == reflect.Ptr {
		dest = dest.Elem()
	}
	q.destination = dest

	criteriaGroups := make([][]criteriaExpression, 0)
	for _, filter := range q.filters {
		criteriaGroup := make([]criteriaExpression, 0)
		for fieldName, value := range filter {
			fieldParts := strings.Split(fieldName, ".")
			switch len(fieldParts) {
			case 1:
				criteriaGroup = append(criteriaGroup, func(datum reflect.Value) bool {
					field := q.model.Fields().GetByName(fieldName)

					return reflect.DeepEqual(value, datum.FieldByIndex(field.Reflection().Index).Interface())
				})
			default:
				panic("indirect fields not implemented")
			}
		}
		criteriaGroups = append(criteriaGroups, criteriaGroup)
	}

	itr := q.txn.tx.GetIterator(make([]byte, 0), false, false)
	items := make([]reflect.Value, 0)
	reader := newDatumReader(q.model)
	for itr.Seek([]byte{datumKeyPrefix}); itr.ValidForPrefix([]byte{datumKeyPrefix}); itr.Next() {
		item := itr.Item()
		key, value, err := make([]byte, 0), make([]byte, 0), error(nil)
		key = item.KeyCopy(key)
		value, err = item.ValueCopy(value)
		if err != nil {
			return err
		}

		if result, err := reader.Read(key, value); err != nil {
			return err
		} else if q.meetsCriteria(result, criteriaGroups) {
			items = append(items, result)
		}
	}

	return nil
}

func (q *Query) meetsCriteria(item reflect.Value, criteria [][]criteriaExpression) bool {
	meetsCriteria := false
	for _, criteriaGroup := range criteria {
		meetsCriteriaGroup := true
		for _, criteria := range criteriaGroup {
			if !criteria(item) {
				meetsCriteriaGroup = false
				break
			}
		}

		if meetsCriteriaGroup {
			meetsCriteria = true
			break
		}
	}

	return meetsCriteria
}
