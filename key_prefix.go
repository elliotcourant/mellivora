package mellivora

type keyPrefix = byte

const (
	datumKeyPrefix = keyPrefix(iota + 1)
	uniqueKeyPrefix
	indexKeyPrefix
)
