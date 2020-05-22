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

	"github.com/dappledger/AnnChain/bcstore/types"
)

func ctxIsDone(ctx context.Context) bool {
	select {
	default:
	case <-ctx.Done():
		return true
	}
	return false
}

type IStore interface {
	Get(context.Context, types.Key) (types.Value, error)
	Set(context.Context, types.KeyValue) error
	Delete(context.Context, types.Key) error
	Close() error
	NewBatch() Batch
	Iterator(prefix types.Key) (Iterator, error)
	IsExist(ctx context.Context, keys types.Key) (bool, error)
}

type Batch interface {
	Set(context.Context, types.KeyValue)
	Delete(context.Context, types.Key)
	Write(context.Context) error
}

type Iterator interface {
	Next(context.Context) bool
	Prev(context.Context) bool
	Key(context.Context) types.Key
	Value(context.Context) types.Value
	Seek(context.Context, types.Key) bool
	Error() error
}
