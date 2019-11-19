package mellivora

type IndexItem map[uint32]interface{}

type IndexIterator interface {
	Seek(values ...interface{})
	Next()
	Valid() bool

	Item() (key []interface{}, value IndexItem)
}
