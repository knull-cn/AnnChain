package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/dappledger/AnnChain/bcstore/proto"
	"github.com/dappledger/AnnChain/bcstore/store"
	"github.com/dappledger/AnnChain/bcstore/types"
)

type Server interface {
	Set(ctx context.Context, kv types.KeyValue) error
	Get(ctx context.Context, key types.Key) (types.KeyValue, error)
	Delete(ctx context.Context, key types.Key) error
	GetByPrefix(ctx context.Context, prefix, lastKey types.Key, limit int64) ([]types.KeyValue, error)
	Batch(ctx context.Context, dels []types.Key, sets []types.KeyValue) error
	Has(ctx context.Context, key types.Key) (bool, error)
	GetProperty(ctx context.Context, property proto.PropertyType) (string, error)
}

func NewServer(dbpath string, dbtype types.StoreType) (Server, error) {
	var is store.IStore
	var err error
	switch dbtype {
	case types.ST_CLevelDB:
		is, err = store.NewCLevelDB(dbpath)
	case types.ST_GOLevelDB:
		is, err = store.NewGoLevelDB(dbpath)
	case types.ST_MemDB:
		is, err = nil, types.SErrNotSupport
	default:
		is, err = nil, types.SErrStoreType
	}
	if err != nil {
		return nil, err
	}
	return &server{
		dbtype,
		is,
	}, err
}

type server struct {
	dbtype types.StoreType
	is     store.IStore
}

func (ls *server) Set(ctx context.Context, kv types.KeyValue) error {
	return ls.is.Set(ctx, kv)
}
func (ls *server) Get(ctx context.Context, key types.Key) (kv types.KeyValue, err error) {
	v, err := ls.is.Get(ctx, key)
	if err != nil {
		return
	}
	kv.Value = v
	kv.Key = key
	return
}

func (ls *server) Delete(ctx context.Context, key types.Key) (err error) {
	return ls.is.Delete(ctx, key)
}

//limit:
func (ls *server) GetByPrefix(ctx context.Context, prefix, lastKey types.Key, limit int64) ([]types.KeyValue, error) {
	if limit <= 0 {
		limit = 20 //default 20;
	}
	itr, err := ls.is.Iterator(prefix)
	if err != nil {
		return nil, err
	}
	if len(lastKey) > 0{
		itr.Seek(ctx, lastKey)
	}
	itr.Next(ctx)
	var kvs []types.KeyValue
	for itr.Next(ctx) {
		kv := types.KeyValue{
			Key:   itr.Key(ctx),
			Value: itr.Value(ctx),
		}
		kvs = append(kvs, kv)
		limit--
		if limit == 0 {
			break
		}
	}
	return kvs, itr.Error()
}
func (ls *server) Batch(ctx context.Context, dels []types.Key, sets []types.KeyValue) error {
	b := ls.is.NewBatch()
	for i := 0; i < len(sets); i++ {
		b.Set(ctx, sets[i])
	}
	for i := 0; i < len(dels); i++ {
		b.Delete(ctx, dels[i])
	}
	return b.Write(ctx)
}
func (ls *server) Has(ctx context.Context, key types.Key) (bool, error) {
	return ls.is.IsExist(ctx, key)
}
func (ls *server) GetProperty(ctx context.Context, property proto.PropertyType) (string, error) {
	if ls.dbtype != types.ST_GOLevelDB {
		return "", fmt.Errorf("db was not GOLevelDB,not support 'GetProperty'")
	}
	var b strings.Builder
	b.WriteString("leveldb.")
	switch property {
	case proto.PropertyType_stats:
		b.WriteString("stats")
	case proto.PropertyType_writedelay:
		b.WriteString("writedelay")
	case proto.PropertyType_iostats:
		b.WriteString("iostats")
	default:
		return "", fmt.Errorf("PropertyType(%s) not support!", property.String())
	}
	db := ls.is.(*store.GoLevelDB)
	return db.GetProperty(b.String())
}
