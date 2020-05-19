package server

import (
	"context"

	"github.com/dappledger/AnnChain/kvstore"
	"github.com/dappledger/AnnChain/kvstore/store"
)

type StoreServer struct {
	dbobj store.IStore
}

func (kvs *StoreServer) StartServer(dbname, dbpath string, dbtype kvstore.StoreType) (err error) {
	kvs.dbobj, err = store.OpenStore(dbname, dbpath, dbtype)
	return err
}

func (kvs *StoreServer) GetList(ctx context.Context, keys []kvstore.KVKey) (ret []kvstore.KVKeyValue, err error) {
	if len(keys) == 0 {
		return nil, kvstore.SErrArg
	}
	if len(keys) > 1 {
		//iterator;
		err = kvstore.SErrNotSupport
	} else {
		v, err := kvs.dbobj.Get(ctx, keys[0])
		if err != nil {
			return nil, err
		}
		kv := kvstore.KVKeyValue{keys[0], v}
		ret = append(ret, kv)
	}
	return ret, err
}
func (kvs *StoreServer) SetList(ctx context.Context, keyvalues []kvstore.KVKeyValue) error {
	if len(keyvalues) == 0 {
		return nil
	}
	if len(keyvalues) == 1 {
		kv := keyvalues[0]
		return kvs.dbobj.Set(ctx, kv)
	}
	//iterator;
	return kvstore.SErrNotSupport
}
func (kvs *StoreServer) DeleteList(ctx context.Context, keys []kvstore.KVKey) error {
	if len(keys) == 0 {
		return nil
	}
	if len(keys) == 1 {
		key := keys[0]
		return kvs.dbobj.Delete(ctx, key)
	}
	//iterator;
	return kvstore.SErrNotSupport
}
