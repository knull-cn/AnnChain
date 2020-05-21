package client

import (
    `flag`

    `github.com/dappledger/AnnChain/bcstore/proto`
)

var (
    rpcaddr = flag.String("rpcaddr","127.0.0.1:9876","address of rpc")
)

func main(){
    proto.NewBCStoreClient()
}