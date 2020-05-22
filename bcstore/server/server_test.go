package server

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dappledger/AnnChain/bcstore/types"
)

var (
	dirs      []string
	keyPrefix = "key_"
)

func TestMain(m *testing.M) {
	r := m.Run()
	finish()
	os.Exit(r)
}

func finish() {
	for _, dir := range dirs {
		err := os.RemoveAll(dir)
		if err != nil {
			fmt.Println("remove(", dir, ")failed :", err, ".")
		}
	}
}

func dbname() string {
	dirname := os.TempDir()
	dirname = path.Join(os.TempDir(), fmt.Sprintf("test_leveldb_%d", time.Now().UnixNano()))
	dirs = append(dirs, dirname)
	return dirname
}

func testGetByPrefix(t *testing.T, db Server, sets []types.KeyValue, nlen int) {
	kvs, err := db.GetByPrefix(context.TODO(), []byte(keyPrefix), nil, int64(nlen))
	assert.Nil(t, err)
	if nlen == len(sets) {
		assert.Equal(t, nlen, len(kvs), fmt.Sprintf("need(%d) buf(%d)", nlen, len(kvs)))
	}
	if nlen == 0 {
		assert.Equalf(t, types.DefPageLimit, len(kvs), fmt.Sprintf("need(%d) buf(%d)", types.DefPageLimit, len(kvs)))
	}
	var filter = make(map[string]struct{}, 10)

	for _, getkv := range kvs {
		_, ok := filter[string(getkv.Key)]
		assert.Falsef(t, ok, "key(%s) has exist", string(getkv.Key))
		//find value;
		for _, kv := range sets {
			if bytes.Equal(getkv.Key, kv.Key) {
				assert.Equalf(t, string(getkv.Value), string(kv.Value), string(kv.Key))
				filter[string(getkv.Key)] = struct{}{}
				break
			}
		}
	}
}

func kv2str(kv *types.KeyValue) string {
	return fmt.Sprintf("%s:%s", string(kv.Key), string(kv.Value))
}

func kvs2str(kvs []types.KeyValue) string {
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < len(kvs); i++ {
		b.WriteString(kv2str(&kvs[i]))
		b.WriteByte(';')
		b.WriteByte(' ')
	}
	b.WriteByte(']')
	b.WriteByte('\n')
	return b.String()
}

func testStoreServer(t *testing.T, dbtype types.StoreType) {
	server, err := NewServer(dbname(), dbtype)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	var times = 32
	ndels := times / 4
	var sets = make([]types.KeyValue, times)
	for i := 0; i < times; i++ {
		sets[i].Key = []byte(fmt.Sprintf("%s%d", keyPrefix, i))
		sets[i].Value = []byte(fmt.Sprintf("value_%d", i))
	}
	//set;
	for _, kv := range sets {
		err = server.Set(context.TODO(), kv)
		assert.Nil(t, err)
	}
	//
	kvs, err := server.GetByPrefix(context.TODO(), []byte(keyPrefix), sets[times-ndels].Key, 30)
	fmt.Printf("lastkey=%s;\nkvs=%v;\nsets=%v.\n", string(sets[times-ndels].Key), kvs2str(kvs), kvs2str(sets[20:]))
	//get;
	for _, kv := range sets {
		getkv, err := server.Get(context.TODO(), kv.Key)
		assert.Nil(t, err)
		assert.Equalf(t, string(kv.Value), string(getkv.Value), string(kv.Key))
	}
	//get by prefix;
	testGetByPrefix(t, server, sets, len(sets)/2)
	testGetByPrefix(t, server, sets, len(sets))
	testGetByPrefix(t, server, sets, 0)
	//deletes;
	dels := sets[:ndels]
	for _, kv := range dels {
		err = server.Delete(context.TODO(), kv.Key)
		assert.Nil(t, err)
	}
	sets = sets[ndels:]
	testGetByPrefix(t, server, sets, len(sets))
}

func TestStoreServer_GoLevelDB(t *testing.T) {
	testStoreServer(t, types.ST_GOLevelDB)
}
