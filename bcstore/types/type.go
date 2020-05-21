package types

import (
	"errors"
)

type Key []byte
type Value []byte

type KeyValue struct {
	Key   Key
	Value Value
}

type StoreType string

const (
	ST_CLevelDB  StoreType = "cleveldb"
	ST_GOLevelDB           = "goleveldb"
	ST_MemDB               = "memdb"
)

//error of server;
var (
	SErrArg        = errors.New("argument format error")
	SErrNotSupport = errors.New("service was not support now")
	SErrStoreType  = errors.New("store type error")
)

var (

)

const (
	DefPageLimit = 20
)