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
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dappledger/AnnChain/bcstore/types"
	gcmn "github.com/dappledger/AnnChain/gemmill/modules/go-common"
)

var (
	keyPrefix = "key_"
)

func testIterator(t *testing.T, db *GoLevelDB, sets []types.KeyValue) {
	itr, err := db.Iterator([]byte(keyPrefix))
	assert.Nil(t, err)
	count := 0
	var filter = make(map[string]struct{}, 10)
	for itr.Next(context.TODO()) {
		key := itr.Key(context.TODO())
		value := itr.Value(context.TODO())
		_, ok := filter[string(key)]
		assert.Falsef(t, ok, "key(%s) has exist", string(key))
		//find value;
		for _, kv := range sets {
			if bytes.Equal(key, kv.Key) {
				assert.Equalf(t, value, kv.Value, string(kv.Key))
				filter[string(key)] = struct{}{}
				break
			}
		}
		count++
	}
	assert.Nil(t, itr.Error())
	assert.Equalf(t, count, len(sets), fmt.Sprintf("len(sets)=(%d) but count=(%d)", len(sets), count))
}

func TestGOLeveldb(t *testing.T) {
	db, err := NewGoLevelDB(dbname())
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	var times = 32
	ndels := times / 4
	var sets = make([]types.KeyValue, times)
	for i := 0; i < times; i++ {
		sets[i].Key = []byte(fmt.Sprintf("%s%d", keyPrefix, i))
		sets[i].Value = RandBytes(32)
	}
	//set;
	for _, kv := range sets {
		err = db.Set(context.TODO(), kv)
		assert.Nil(t, err)
	}
	//get;
	for _, kv := range sets {
		getkv, err := db.Get(context.TODO(), kv.Key)
		assert.Nil(t, err)
		assert.Equalf(t, kv.Value, getkv, string(kv.Key))
	}
	//get by prefix;
	testIterator(t, db, sets)
	//deletes;
	dels := sets[:ndels]
	for _, kv := range dels {
		err = db.Delete(context.TODO(), kv.Key)
		assert.Nil(t, err)
	}
	sets = sets[ndels:]
	testIterator(t, db, sets)
}

func BenchmarkRandomReadsWrites(b *testing.B) {
	b.StopTimer()

	numItems := int64(1000000)
	internal := map[int64]int64{}
	for i := 0; i < int(numItems); i++ {
		internal[int64(i)] = int64(0)
	}
	db, err := NewGoLevelDB(dbname())
	if err != nil {
		b.Fatal(err.Error())
		return
	}

	fmt.Println("ok, starting")
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		// Write something
		{
			idx := RandInt64(numItems)
			internal[idx] += 1
			val := internal[idx]
			idxBytes := int642Bytes(int64(idx))
			valBytes := int642Bytes(int64(val))
			//fmt.Printf("Set %X -> %X\n", idxBytes, valBytes)
			err := db.Set(
				context.TODO(),
				types.KeyValue{
					Key:   idxBytes,
					Value: valBytes,
				},
			)
			assert.Nil(b, err)
		}
		// Read something
		{
			idx := (int64(gcmn.RandInt()) % numItems)
			val := internal[idx]
			idxBytes := int642Bytes(idx)
			valBytes, err := db.Get(context.TODO(), idxBytes)
			assert.Nil(b, err)
			//fmt.Printf("Get %X -> %X\n", idxBytes, valBytes)
			if val == 0 {
				if !bytes.Equal(valBytes, nil) {
					b.Errorf("Expected %v for %v, got %X",
						nil, idx, valBytes)
					break
				}
			} else {
				if len(valBytes) != 8 {
					b.Errorf("Expected length 8 for %v, got %X",
						idx, valBytes)
					break
				}
				valGot := bytes2Int64(valBytes)
				if val != valGot {
					b.Errorf("Expected %v for %v, got %v",
						val, idx, valGot)
					break
				}
			}
		}
	}

	db.Close()
}

func int642Bytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func bytes2Int64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
