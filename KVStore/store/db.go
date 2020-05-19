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
	`context`

	`github.com/dappledger/AnnChain/kvstore`
)

type IStore interface {
	Get(context.Context, kvstore.KVKey) (kvstore.KVValue, error)
	Set(context.Context, kvstore.KVKeyValue) error
	Delete(context.Context, kvstore.KVKey) error
	Close() error
	NewBatch() Batch
	Iterator(beg,end kvstore.KVKey) (Iterator,error)
}

type Batch interface {
	Set(context.Context, kvstore.KVKeyValue)
	Delete(context.Context, kvstore.KVKey)
	Write(context.Context) error
}

type Iterator interface {
	Next(context.Context) bool
	Prev(context.Context) bool
	Key(context.Context) kvstore.KVKey
	Value(context.Context) kvstore.KVValue
	Error() error
}



//-----------------------------------------------------------------------------



//	registerDBCreator(LevelDBBackendStr, dbCreator, true)
//	registerDBCreator(CLevelDBBackendStr, dbCreator, false)

//go_level_db
//registerDBCreator(LevelDBBackendStr, dbCreator, false)
//registerDBCreator(GoLevelDBBackendStr, dbCreator, false)
func OpenStore(name ,path string,dbtype kvstore.StoreType)(IStore,error){
	switch dbtype {
	case kvstore.ST_CLevelDB:
		return NewCLevelDB(name,path)
	case kvstore.ST_GOLevelDB:
		return NewGoLevelDB(name,path)
	case kvstore.ST_MemDB:
		return nil,kvstore.SErrNotSupport
	default:
		return nil,kvstore.SErrStoreType
	}
}
