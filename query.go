package mellivora

type Query struct {
	model Model
	txn   *Transaction
}

func (q *Query) Where(on func(criteria *Criteria) *Criteria) *Query {
	return q
}

func (q *Query) LeftJoin(modelName, asName string, on func(criteria *Criteria) *Criteria) *Query {

}
