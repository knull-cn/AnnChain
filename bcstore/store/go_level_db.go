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

package store

import (
	"context"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	lvlitr "github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/dappledger/AnnChain/bcstore/types"
)

type GoLevelDB struct {
	db *leveldb.DB
}

func NewGoLevelDB(dbPath string) (*GoLevelDB, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	database := &GoLevelDB{db: db}
	return database, nil
}

func (db *GoLevelDB) Get(ctx context.Context, key types.Key) (types.Value, error) {
	res, err := db.db.Get(key, nil)
	if err != nil {
		if err != errors.ErrNotFound {
			return nil, err
		}
	}
	return res, nil
}

func (db *GoLevelDB) Set(ctx context.Context, kv types.KeyValue) error {
	return db.db.Put(kv.Key, kv.Value, nil)
}

func (db *GoLevelDB) SetSync(kv types.KeyValue) error {
	return db.db.Put(kv.Key, kv.Value, &opt.WriteOptions{Sync: true})
}

func (db *GoLevelDB) Delete(ctx context.Context, key types.Key) error {
	return db.db.Delete(key, nil)
}

func (db *GoLevelDB) DeleteSync(ctx context.Context, key types.Key) error {
	return db.db.Delete(key, &opt.WriteOptions{Sync: true})
}

func (db *GoLevelDB) Close() error {
	return db.db.Close()
}

func (db *GoLevelDB) GetProperty(name string) (string, error) {
	return db.db.GetProperty(name)
}

func (db *GoLevelDB) IsExist(ctx context.Context, key types.Key) (bool, error) {
	return db.db.Has(key, nil)
}

func (db *GoLevelDB) Print() {
	iter := db.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("[%X]:\t[%X]\n", key, value)
	}
}

func (db *GoLevelDB) Iterator(prefix types.Key) (Iterator, error) {
	ii := db.db.NewIterator(util.BytesPrefix(prefix), nil)
	return &goLevelIterator{
		prefix,
		ii,
	}, nil
}

func (db *GoLevelDB) NewBatch() Batch {
	batch := new(leveldb.Batch)
	return &goLevelDBBatch{
		db,
		batch,
	}
}

type goLevelIterator struct {
	prefix types.Key
	ii     lvlitr.Iterator
}

func (ic *goLevelIterator) Next(context.Context) bool {
	return ic.ii.Next()
}

func (ic *goLevelIterator) Prev(context.Context) bool {
	return ic.ii.Prev()
}
func (ic *goLevelIterator) Key(context.Context) types.Key {
	return ic.ii.Key()
}
func (ic *goLevelIterator) Value(context.Context) types.Value {
	return ic.ii.Value()
}

func (ic *goLevelIterator) Seek(ctx context.Context, key types.Key) bool {
	return ic.ii.Seek(key)
}

func (ic *goLevelIterator) Error() error {
	return ic.ii.Error()
}

//--------------------------------------------------------------------------------

type goLevelDBBatch struct {
	db    *GoLevelDB
	batch *leveldb.Batch
}

func (mBatch *goLevelDBBatch) Set(ctx context.Context, kv types.KeyValue) {
	mBatch.batch.Put(kv.Key, kv.Value)
}

func (mBatch *goLevelDBBatch) Delete(ctx context.Context, key types.Key) {
	mBatch.batch.Delete(key)
}

func (mBatch *goLevelDBBatch) Write(ctx context.Context) error {
	return mBatch.db.db.Write(mBatch.batch, nil)
}
