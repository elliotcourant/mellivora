package mellivora

type (
	StorageEngine interface {
		Begin() (*Transaction, error)
	}

	TransactionEngine interface {
	}
)
