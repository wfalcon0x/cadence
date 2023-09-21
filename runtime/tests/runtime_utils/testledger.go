/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package runtime_utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/onflow/atree"
)

type TestLedger struct {
	StoredValues           map[string][]byte
	OnValueExists          func(owner, key []byte) (exists bool, err error)
	OnGetValue             func(owner, key []byte) (value []byte, err error)
	OnSetValue             func(owner, key, value []byte) (err error)
	OnAllocateStorageIndex func(owner []byte) (atree.StorageIndex, error)
}

var _ atree.Ledger = TestLedger{}

func (s TestLedger) GetValue(owner, key []byte) (value []byte, err error) {
	return s.OnGetValue(owner, key)
}

func (s TestLedger) SetValue(owner, key, value []byte) (err error) {
	return s.OnSetValue(owner, key, value)
}

func (s TestLedger) ValueExists(owner, key []byte) (exists bool, err error) {
	return s.OnValueExists(owner, key)
}

func (s TestLedger) AllocateStorageIndex(owner []byte) (atree.StorageIndex, error) {
	return s.OnAllocateStorageIndex(owner)
}

func (s TestLedger) Dump() {
	for key, data := range s.StoredValues {
		fmt.Printf("%s:\n", strconv.Quote(key))
		fmt.Printf("%s\n", hex.Dump(data))
		println()
	}
}

func NewTestLedger(
	onRead func(owner, key, value []byte),
	onWrite func(owner, key, value []byte),
) TestLedger {

	storageKey := func(owner, key string) string {
		return strings.Join([]string{owner, key}, "|")
	}

	storedValues := map[string][]byte{}

	storageIndices := map[string]uint64{}

	return TestLedger{
		StoredValues: storedValues,
		OnValueExists: func(owner, key []byte) (bool, error) {
			value := storedValues[storageKey(string(owner), string(key))]
			return len(value) > 0, nil
		},
		OnGetValue: func(owner, key []byte) (value []byte, err error) {
			value = storedValues[storageKey(string(owner), string(key))]
			if onRead != nil {
				onRead(owner, key, value)
			}
			return value, nil
		},
		OnSetValue: func(owner, key, value []byte) (err error) {
			storedValues[storageKey(string(owner), string(key))] = value
			if onWrite != nil {
				onWrite(owner, key, value)
			}
			return nil
		},
		OnAllocateStorageIndex: func(owner []byte) (result atree.StorageIndex, err error) {
			index := storageIndices[string(owner)] + 1
			storageIndices[string(owner)] = index
			binary.BigEndian.PutUint64(result[:], index)
			return
		},
	}
}
