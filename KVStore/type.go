package kvstore

import (
	"errors"
)

type KVKey []byte
type KVValue []byte

type KVKeyValue struct {
	Key   KVKey
	Value KVValue
}

type StoreType int

const (
	ST_CLevelDB  StoreType= iota //"cleveldb"
	ST_GOLevelDB        //"goleveldb"
	ST_MemDB            //"memdb"
)

var (
	SErrArg        = errors.New("argument format error")
	SErrNotSupport = errors.New("service was not support now")
	SErrStoreType  = errors.New("store type error")
)
