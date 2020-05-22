package main

import (
	"flag"
	"fmt"
	"net"
	`os`

	"google.golang.org/grpc"

	"github.com/dappledger/AnnChain/bcstore/proto"
	"github.com/dappledger/AnnChain/bcstore/server"
	`github.com/dappledger/AnnChain/bcstore/types`
	"github.com/sirupsen/logrus"
)

var (
	portstr = flag.String("port", "9876", "listener of this service")
	runtime = flag.String("runtime","","database dir for runtime")
	dbtype = flag.Int64("dbtype",0,"database type(default is 0).\n\t0: CLevelDB;\n\t1: GOLevelDB;\n\t2: MemDB.\n")
	isDebug = flag.Bool("debug",false,"is debug model.true is debug")
)

func loginit(){
	logrus.SetOutput(os.Stdout)
	lvl := logrus.InfoLevel
	if *isDebug{
		lvl = logrus.TraceLevel
	}
	logrus.SetLevel(lvl)
}

func main() {
	flag.Parse()
	loginit()
	l, err := net.Listen("tcp", ":"+*portstr)
	if err != nil {
		panic("listen err:" + err.Error())
	}
	//var service = server.NewStoreService(viper.GetString("runtime"),types.StoreType(viper.GetInt64("dbtype")))
	var service = server.NewStoreService(*runtime,types.StoreType(*dbtype))
	s := grpc.NewServer()
	proto.RegisterBCStoreServer(s, service)
	err = s.Serve(l)
	if err != nil {
		fmt.Println("Serve error:", err)
	}
}
