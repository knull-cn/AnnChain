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

//--- +build gcc

package store

import (
	"context"
	"fmt"

	"github.com/jmhodges/levigo"

	"github.com/dappledger/AnnChain/bcstore/types"
)

type CLevelDB struct {
	db     *levigo.DB
	ro     *levigo.ReadOptions
	wo     *levigo.WriteOptions
	woSync *levigo.WriteOptions
}

func NewCLevelDB(dbpath string) (*CLevelDB, error) {

	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(1 << 30))
	opts.SetCreateIfMissing(true)
	db, err := levigo.Open(dbpath, opts)
	if err != nil {
		return nil, err
	}
	ro := levigo.NewReadOptions()
	wo := levigo.NewWriteOptions()
	woSync := levigo.NewWriteOptions()
	woSync.SetSync(true)
	database := &CLevelDB{
		db:     db,
		ro:     ro,
		wo:     wo,
		woSync: woSync,
	}
	return database, nil
}

func ctxIsDone(ctx context.Context) bool {
	select {
	default:
	case <-ctx.Done():
		return true
	}
	return false
}

func (db *CLevelDB) Get(ctx context.Context, key types.Key) (types.Value, error) {
	if ctxIsDone(ctx){
		return nil,ctx.Err()
	}
	return db.db.Get(db.ro, key)
}

func (db *CLevelDB) Set(ctx context.Context, kv types.KeyValue) error {
	if ctxIsDone(ctx){
		return ctx.Err()
	}
	return db.db.Put(db.wo, kv.Key, kv.Value)
}

func (db *CLevelDB) SetSync(ctx context.Context, kv types.KeyValue) error {
	if ctxIsDone(ctx){
		return ctx.Err()
	}
	return db.db.Put(db.woSync, kv.Key, kv.Value)
}

func (db *CLevelDB) Delete(ctx context.Context, key types.Key) error {
	if ctxIsDone(ctx){
		return ctx.Err()
	}
	return db.db.Delete(db.wo, key)
}

func (db *CLevelDB) DeleteSync(ctx context.Context, key types.Key) error {
	if ctxIsDone(ctx){
		return ctx.Err()
	}
	return db.db.Delete(db.woSync, key)
}

func (db *CLevelDB) IsExist(ctx context.Context, key types.Key) (bool, error) {
	if ctxIsDone(ctx){
		return false,ctx.Err()
	}
	v, err := db.db.Get(db.ro, key)
	return (v != nil), err
}

func (db *CLevelDB) DB() *levigo.DB {
	return db.db
}

func (db *CLevelDB) Close() error {
	db.db.Close()
	db.ro.Close()
	db.wo.Close()
	db.woSync.Close()
	return nil
}

func (db *CLevelDB) Print() {
	iter := db.db.NewIterator(db.ro)
	defer iter.Close()
	for iter.Seek(nil); iter.Valid(); iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("[%X]:\t[%X]\n", key, value)
	}
}

func (db *CLevelDB) NewBatch() Batch {
	batch := levigo.NewWriteBatch()
	return &cLevelDBBatch{db, batch}
}

func (db *CLevelDB) Iterator(prefix types.Key) (Iterator, error) {
	itr := db.db.NewIterator(db.ro)
	itr.Seek(prefix)
	return &cLevelIter{
		prefix,
		itr,
	}, nil
}

type cLevelIter struct {
	prefix types.Key
	ii     *levigo.Iterator
}

func (ic *cLevelIter) Next( context.Context) bool {
	ic.ii.Next()
	return ic.ii.GetError() == nil
}

func (ic *cLevelIter) Prev(context.Context) bool {
	ic.ii.Prev()
	return ic.ii.GetError() == nil
}
func (ic *cLevelIter) Key(context.Context) types.Key {
	return ic.ii.Key()
}
func (ic *cLevelIter) Value(context.Context) types.Value {
	return ic.ii.Value()
}

func (ic *cLevelIter) Seek(_ context.Context, key types.Key) bool {
	ic.ii.Seek(key)
	return ic.Error() == nil
}

func (ic *cLevelIter) Error() error {
	return ic.ii.GetError()
}

//--------------------------------------------------------------------------------

type cLevelDBBatch struct {
	db    *CLevelDB
	batch *levigo.WriteBatch
}

func (mBatch *cLevelDBBatch) Set(_ context.Context, kv types.KeyValue) {
	mBatch.batch.Put(kv.Key, kv.Value)
}

func (mBatch *cLevelDBBatch) Delete(_ context.Context, key types.Key) {
	mBatch.batch.Delete(key)
}

func (mBatch *cLevelDBBatch) Write(_ context.Context) error {
	return mBatch.db.db.Write(mBatch.db.wo, mBatch.batch)
}
