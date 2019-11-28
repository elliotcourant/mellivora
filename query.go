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

	return nil
}
