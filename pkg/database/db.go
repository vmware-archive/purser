/*
 * Copyright (c) 2018 VMware Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package database

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// PutInDB stores the key value pair in the db's bucket
// Usage:
// 		myStruct := MyStruct{k1: v1, k2: v2 ...}
//		err := PutInDB(yourDB, bucketName, keyName, myStruct)
func PutInDB(boltDB *bolt.DB, bucket string, key interface{}, value interface{}) error {
	err := boltDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket %s: ", err)
		}

		keyInBytes, errKey := json.Marshal(key)
		if errKey != nil {
			return fmt.Errorf("json marshal %s: ", errKey)
		}

		valueInBytes, errVal := json.Marshal(value)
		if errVal != nil {
			return fmt.Errorf("json marshal %s: ", errVal)
		}

		err = b.Put(keyInBytes, valueInBytes)
		return err
	})
	return err
}

// GetValueFromDB gives the value for corresponding key in the db's bucket
// Usage:
//		myStruct := &MyStruct{}
//		err := GetValueFromDB(yourDB, bucketName, keyName, myStruct)
func GetValueFromDB(boltDB *bolt.DB, bucket string, key interface{}, value interface{}) error {
	err := boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		keyInBytes, errKey := json.Marshal(key)
		if errKey != nil {
			return fmt.Errorf("json marshal %s: ", errKey)
		}

		valueInBytes := b.Get(keyInBytes)
		err := json.Unmarshal(valueInBytes, value)
		return err
	})
	return err
}
