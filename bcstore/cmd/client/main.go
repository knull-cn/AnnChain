package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	bcclient "github.com/dappledger/AnnChain/bcstore/client"
	"github.com/dappledger/AnnChain/bcstore/types"
	"github.com/sirupsen/logrus"
)

var (
	rootCmd = &cobra.Command{
		Use:   "kvcli",
		Short: "client for kvstore",
	}
	//
	dbname = "test"
	dbdir  = ""
	dbtype = types.StoreType(types.ST_GOLevelDB)
)

func kv2str(kv types.KeyValue) string {
	return fmt.Sprintf("%s:%s", string(kv.Key), string(kv.Value))
}

func main() {
	bcclient.FlagsSet(rootCmd)
	rootCmd.PersistentFlags().Set(bcclient.FlagUseRPC, "true")
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(os.Stdout)
	dbstore, err := bcclient.NewClient(dbname, dbdir, dbtype)
	if err != nil {
		fmt.Printf("NewClient(%s) error:%s", dbname, err.Error())
		return
	}
	var cmd string
	var key string
	var value string
	var args []string
	running := true
	stdin := bufio.NewReader(os.Stdin)
	for running {
		args = make([]string, 0, 3)
		fmt.Println("chose your command: 'set $key $value/get $key/getbp $prefix;any other for quit.'")
		fmt.Print("~> ")
		line, _, _ := stdin.ReadLine()
		arrs := strings.Split(string(line), " ")
		for i := 0; i < len(arrs); i++ {
			tmp := strings.TrimSpace(arrs[i])
			if len(tmp) > 0 {
				args = append(args, tmp)
			}
		}
		cmd = args[0]
		fmt.Printf("line=%s;args=%v.\n", line, args)
		switch cmd {
		case "set":
			if len(args) < 3 {
				fmt.Println("command 'set' args need 2,but there is ", len(args)-1)
				continue
			}
			key = args[1]
			value = args[2]
			fmt.Printf("now we do : 'set %s %s'.\n", key, value)
			err = dbstore.Set(types.KeyValue{[]byte(key), []byte(value)})
			if err == nil {
				fmt.Println("~> OK.")
			} else {
				fmt.Println("~> Error :", err)
			}
		case "get":
			if len(args) < 2 {
				fmt.Println("command 'get' args need `1`,but there is no arg.")
				continue
			}
			key = args[1]
			fmt.Printf("now we do : 'get %s'.\n", key)
			kv, err := dbstore.Get(types.Key(key))
			if err == nil {
				fmt.Printf("~> OK. %s.\n", kv2str(kv))
			} else {
				fmt.Println("~> Error :", err)
			}
		case "getbp":
			if len(args) < 2 {
				fmt.Println("command 'getbp' args need `1`,but there is no arg.")
				continue
			}
			key = args[1]
			fmt.Printf("now we do : 'get %s'.\n", key)

			kvs, err := dbstore.GetByPrefix(types.Key(key), types.Key{}, 20)
			if err == nil {
				fmt.Printf("~> Ok. getByPrefix('%s*'):\n", key)
				for _, kv := range kvs {
					fmt.Printf("\t%s;", kv2str(kv))
				}
				fmt.Println("")
			} else {
				fmt.Println("~> Error :", err)
			}
		default:
			fmt.Printf("command('%s') do exit!.\n", cmd)
			running = false
		}
	}
}
