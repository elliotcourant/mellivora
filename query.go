package mellivora

import (
	"fmt"
	"reflect"
	"strings"
	"time"
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
	start := time.Now()
	defer q.txn.db.logger.Tracef("select %T took %s", destination, time.Since(start))

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

					// TODO (elliotcourant) build a better comparision system.
					return fmt.Sprint(value) == fmt.Sprint(datum.FieldByIndex(field.Reflection().Index).Interface())
				})
			default:
				panic("indirect fields not implemented")
			}
		}
		criteriaGroups = append(criteriaGroups, criteriaGroup)
	}

	itr := q.txn.iterator()
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

	return q.scanResults(items)
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

func (q *Query) scanResults(items []reflect.Value) error {
	switch q.destination.Kind() {
	case reflect.Struct:
		q.destination.Set(items[0])
	case reflect.Array, reflect.Slice:
		q.destination.Set(reflect.MakeSlice(q.destination.Type(), 0, q.destination.Cap()))
		q.destination.Set(reflect.Append(q.destination, items...))
	default:
		return fmt.Errorf("cannot scan results to %T", q.destination)
	}

	return nil
}
