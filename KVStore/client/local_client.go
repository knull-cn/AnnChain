package client

import (
    "context"

    "github.com/dappledger/AnnChain/kvstore"
    `github.com/dappledger/AnnChain/kvstore/server`
)

type LocalClient struct {
	srv server.StoreServer
	dbname string
	dbpath string
	dbtype string
}

func NewLocalClient(dbname,dbpath string,dbtype kvstore.StoreType)(Client,error){
    var lc LocalClient
    err:= lc.srv.StartServer(dbname,dbpath,dbtype)
    return &lc,err
}

func (kc *LocalClient) ConnectServer(addrs []string) error {
	return nil
}
func (kc *LocalClient) GetList(ctx context.Context, keys []kvstore.KVKey) ([]kvstore.KVKeyValue, error) {
    return kc.srv.GetList(ctx,keys)
}
func (kc *LocalClient) SetList(ctx context.Context, kvs []kvstore.KVKeyValue) error{
    return kc.srv.SetList(ctx,kvs)
}
func (kc *LocalClient) DeleteList(ctx context.Context, keys []kvstore.KVKey) error{
    return kc.srv.DeleteList(ctx,keys)
}
