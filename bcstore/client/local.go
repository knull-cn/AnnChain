package client

import (
	"context"

	"github.com/dappledger/AnnChain/bcstore/proto"
	"github.com/dappledger/AnnChain/bcstore/server"
	"github.com/dappledger/AnnChain/bcstore/types"
)

type DBStore interface {
	Set(key types.KeyValue) error
	Get(key types.Key) (types.KeyValue, error)
	Delete(key types.Key) error
	GetByPrefix(prefix, lastKey types.Key, limit int64) ([]types.KeyValue, error)
	Batch(dels []types.Key, sets []types.KeyValue) error
	Has(key types.Key) (bool, error)
	GetProperty(property proto.PropertyType) (string, error)
}

type LocalClient struct {
	srv server.Server
}

func NewLocalClient(fullpath string, dbtype types.StoreType) (ds DBStore, err error) {
	var lc LocalClient
	lc.srv, err = server.NewServer(fullpath, dbtype)
	if err != nil {
		return nil, err
	}
	return &lc, nil
}

func (kc *LocalClient) ConnectServer(addrs []string) error {
	return nil
}
func (kc *LocalClient) Set(kv types.KeyValue) error {
	return kc.srv.Set(context.TODO(), kv)
}
func (kc *LocalClient) Get(key types.Key) (types.KeyValue, error) {
	return kc.srv.Get(context.TODO(), key)
}

func (kc *LocalClient) Delete(key types.Key) error {
	return kc.srv.Delete(context.TODO(), key)
}
func (kc *LocalClient) GetByPrefix(prefix, lastKey types.Key, limit int64) ([]types.KeyValue, error) {
	return kc.srv.GetByPrefix(context.TODO(), prefix, lastKey, limit)
}

func (kc *LocalClient) Has(key types.Key) (bool, error) {
	return kc.srv.Has(context.TODO(), key)
}
func (kc *LocalClient) Batch(dels []types.Key, sets []types.KeyValue) error {
	return kc.srv.Batch(context.TODO(), dels, sets)
}
func (kc *LocalClient) GetProperty(name proto.PropertyType) (string, error) {
	return kc.srv.GetProperty(context.TODO(), name)
}
