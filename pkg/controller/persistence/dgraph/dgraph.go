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

package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

// Dgraph variables
var (
	Client     *dgo.Dgraph
	Connection *grpc.ClientConn
)

func init() {
	err := Open("127.0.0.1:9080")
	if err != nil {
		fmt.Println("Error while opening connection to Dgraph ", err)
	}

	err = CreateSchema()
	if err != nil {
		fmt.Println("Error while creating schema ", err)
	}
}

// Open creates and establishes a new Dgraph connection
func Open(url string) error {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return err
	}

	Connection = conn
	dc := api.NewDgraphClient(Connection)
	Client = dgo.NewDgraphClient(dc)

	return nil
}

// Close terminates the Dgraph connection
func Close() {
	err := Connection.Close()
	if err != nil {
		fmt.Println("Error closing connection to Dgraph ", err)
	}
}

// CreateSchema sets the Dgraph schema
func CreateSchema() error {
	op := &api.Operation{}
	op.Schema = `
		name: string @index(term) .
		xid:  string @index(term) .
		isService: bool .
		isPod: bool .
		isContainer: bool .
	`
	ctx := context.Background()
	err := Client.Alter(ctx, op)

	return err
}

// GetUID returns the UID of the node in the Dgraph
func GetUID(dg *dgo.Dgraph, id string, nodeType string) string {

	query := `query Me($id:string, $nodeType:string) {
		getUid(func: eq(xid, $id)) @filter(has(` + nodeType + `)) {
			uid
		}
	}`

	ctx := context.Background()
	variables := make(map[string]string)
	variables["$nodeType"] = nodeType
	variables["$id"] = id

	resp, err := dg.NewReadOnlyTxn().QueryWithVars(ctx, query, variables)
	if err != nil {
		log.Printf("failed to fetch UID from Dgraph %v", err)
		return ""
	}
	return unmarshalDgraphResponse(resp, id)
}

// MutateNode mutates a Dgraph transaction
func MutateNode(dg *dgo.Dgraph, n []byte) (*api.Assigned, error) {
	mu := &api.Mutation{
		CommitNow: true,
	}

	mu.SetJson = n
	ctx := context.Background()
	return dg.NewTxn().Mutate(ctx, mu)
}

func unmarshalDgraphResponse(resp *api.Response, id string) string {
	type Root struct {
		IDs []ID `json:"getUid"`
	}

	var r Root
	err := json.Unmarshal(resp.Json, &r)
	if err != nil {
		log.Printf("failed to marshal Dgraph response %v", err)
		return ""
	}

	if len(r.IDs) == 0 {
		log.Printf("id %s is not in dgraph", id)
		return ""
	}

	return r.IDs[0].UID
}
