package client

import (
    "context"

    `github.com/dappledger/AnnChain/kvstore`
)

//this is call by annchain process and so on;
type Client interface {
	ConnectServer(addrs []string) error
	GetList(context.Context, []kvstore.KVKey) ([]kvstore.KVKeyValue, error)
	SetList(context.Context, []kvstore.KVKeyValue) error
	DeleteList(context.Context, []kvstore.KVKey) error
}
