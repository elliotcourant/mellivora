package mellivora

type Ex map[string]interface{}

type Query struct {
	model   Model
	txn     *Transaction
	filters []Ex

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
	return nil
}
