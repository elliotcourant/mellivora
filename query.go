package mellivora

type Query struct {
	model interface{}
	txn   *Transaction
}
