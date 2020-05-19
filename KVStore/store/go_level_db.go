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
	"path"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	lvlitr "github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/dappledger/AnnChain/kvstore"
)

type GoLevelDB struct {
	db *leveldb.DB
}

func NewGoLevelDB(name string, dir string) (*GoLevelDB, error) {
	dbPath := path.Join(dir, name+".db")
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	database := &GoLevelDB{db: db}
	return database, nil
}

func (db *GoLevelDB) Get(ctx context.Context, key kvstore.KVKey) (kvstore.KVValue, error) {
	res, err := db.db.Get(key, nil)
	if err != nil {
		if err != errors.ErrNotFound {
			return nil, err
		}
	}
	return res, nil
}

func (db *GoLevelDB) Set(ctx context.Context, kv kvstore.KVKeyValue) error {
	return db.db.Put(kv.Key, kv.Value, nil)
}

func (db *GoLevelDB) SetSync(kv kvstore.KVKeyValue) error {
	return db.db.Put(kv.Key, kv.Value, &opt.WriteOptions{Sync: true})
}

func (db *GoLevelDB) Delete(ctx context.Context, key kvstore.KVKey) error {
	return db.db.Delete(key, nil)
}

func (db *GoLevelDB) DeleteSync(ctx context.Context, key kvstore.KVKey) error {
	return db.db.Delete(key, &opt.WriteOptions{Sync: true})
}

func (db *GoLevelDB) Close() error {
	return db.db.Close()
}

func (db *GoLevelDB) Print() {
	iter := db.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("[%X]:\t[%X]\n", key, value)
	}
}

func (db *GoLevelDB) Iterator(beg, end kvstore.KVKey) (Iterator, error) {
	ii := db.db.NewIterator(&util.Range{beg, end}, nil)
	return &goLevelIterator{
		beg,
		end,
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
	start, end kvstore.KVKey
	ii         lvlitr.Iterator
}

func (ic *goLevelIterator) Next(context.Context) bool {
	return ic.ii.Next()
}

func (ic *goLevelIterator) Prev(context.Context) bool {
	return ic.ii.Prev()
}
func (ic *goLevelIterator) Key(context.Context) kvstore.KVKey {
	return ic.ii.Key()
}
func (ic *goLevelIterator) Value(context.Context) kvstore.KVValue {
	return ic.ii.Value()
}

func (ic *goLevelIterator) Error() error {
	return ic.ii.Error()
}

//--------------------------------------------------------------------------------

type goLevelDBBatch struct {
	db    *GoLevelDB
	batch *leveldb.Batch
}

func (mBatch *goLevelDBBatch) Set(ctx context.Context, kv kvstore.KVKeyValue) {
	mBatch.batch.Put(kv.Key, kv.Value)
}

func (mBatch *goLevelDBBatch) Delete(ctx context.Context, key kvstore.KVKey) {
	mBatch.batch.Delete(key)
}

func (mBatch *goLevelDBBatch) Write(ctx context.Context) error {
	return mBatch.db.db.Write(mBatch.batch, nil)
}
