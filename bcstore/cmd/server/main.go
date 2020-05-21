package server

import (
	"flag"
	"fmt"
	"net"

	`github.com/spf13/viper`
	"google.golang.org/grpc"

	"github.com/dappledger/AnnChain/bcstore/proto"
	"github.com/dappledger/AnnChain/bcstore/server"
	`github.com/dappledger/AnnChain/bcstore/types`
)

var (
	portstr = flag.String("port", "9876", "listener of this service")
)

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", ":"+*portstr)
	if err != nil {
		panic("listen err:" + err.Error())
	}
	var service = server.NewStoreService(viper.GetString("runtime"),types.StoreType(viper.GetInt64("dbtype")))
	s := grpc.NewServer()
	proto.RegisterBCStoreServer(s, service)
	err = s.Serve(l)
	if err != nil {
		fmt.Println("Serve error:", err)
	}
}
