// Copyright 2017 ZhongAn Information Technology Services Co.,Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	kvclient "github.com/dappledger/AnnChain/bcstore/client"
	"github.com/dappledger/AnnChain/bcstore/types"
	gcmn "github.com/dappledger/AnnChain/gemmill/modules/go-common"
)

type DB interface {
	Get([]byte) []byte
	Set([]byte, []byte)
	SetSync([]byte, []byte)
	Delete([]byte)
	DeleteSync([]byte)
	Close()
	NewBatch() Batch

	// For debugging
	Print()
	Iterator() Iterator
}

type Batch interface {
	Set(key, value []byte)
	Delete(key []byte)
	Write()
}

type Iterator interface {
	Next() bool
	Key() []byte
	Value() []byte
}

//-----------------------------------------------------------------------------

const (
	LevelDBBackendStr   = "leveldb" // legacy, defaults to goleveldb.
	CLevelDBBackendStr  = "cleveldb"
	GoLevelDBBackendStr = "goleveldb"
	MemDBBackendStr     = "memdb"
)

type dbCreator func(name string, dir string) (DB, error)

var backends = map[string]dbCreator{}

func registerDBCreator(backend string, creator dbCreator, force bool) {
	_, ok := backends[backend]
	if !force && ok {
		return
	}
	backends[backend] = creator
}

func NewDB(name string, backend string, dir string) DB {
	ds ,err := NewDBErr(name,backend,dir)
	if err != nil {
		gcmn.PanicSanity(gcmn.Fmt("Error initializing DB: %v", err))
	}
	return ds
}

func NewDBErr(name string, backend string, dir string) (DB,error) {
	var ds kvclient.DBStore
	var err error
	switch backend {
	case LevelDBBackendStr, GoLevelDBBackendStr:
		ds, err = kvclient.NewClient(name,dir, types.ST_GOLevelDB)
	case CLevelDBBackendStr:
		ds, err = kvclient.NewClient(name,dir, types.ST_CLevelDB)
	case MemDBBackendStr:
		ds, err = kvclient.NewClient(name,dir, types.ST_MemDB)
	default:
		err = types.SErrNotSupport
	}
	return &KVClient{ds},err
}

type KVClient struct {
	ds kvclient.DBStore
}

func (kvc *KVClient) Set(k, v []byte) {
	err := kvc.ds.Set(types.KeyValue{k, v})
	if err != nil {
		gcmn.PanicCrisis(err)
	}
}
func (kvc *KVClient) SetSync(k []byte, v []byte) {
	kvc.Set(k, v)
}
func (kvc *KVClient) Delete(key []byte) {
	err := kvc.ds.Delete(key)
	if err != nil {
		gcmn.PanicCrisis(err)
	}
}
func (kvc *KVClient) DeleteSync(key []byte) {
	err := kvc.ds.Delete(key)
	if err != nil {
		gcmn.PanicCrisis(err)
	}
}
func (kvc *KVClient) Close() {

}
func (kvc *KVClient) NewBatch() Batch {
	return &KVBatch{
		kvc: kvc,
	}
}
func (kvc *KVClient) Print() {

}
func (kvc *KVClient) Iterator() Iterator {
	return nil
}

func (kvc *KVClient) Get(key []byte) []byte {
	kv, err := kvc.ds.Get(key)
	if err != nil {
		gcmn.PanicCrisis(err)
	}
	return kv.Value
}

type KVBatch struct {
	kvc     *KVClient
	setkvs  []types.KeyValue
	delkeys []types.Key
}

func (kvb *KVBatch) Set(key, value []byte) {
	kvb.setkvs = append(kvb.setkvs, types.KeyValue{key, value})
}
func (kvb *KVBatch) Delete(key []byte) {
	kvb.delkeys = append(kvb.delkeys, key)
}
func (kvb *KVBatch) Write() {
	kvb.kvc.ds.Batch(kvb.delkeys, kvb.setkvs)
}
