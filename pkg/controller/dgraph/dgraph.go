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

	log "github.com/Sirupsen/logrus"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/vmware/purser/pkg/controller/utils"
	"google.golang.org/grpc"
)

// mutation types
const (
	CREATE = "create"
	UPDATE = "update"
	DELETE = "delete"
)

// Dgraph variables
var (
	client     *dgo.Dgraph
	connection *grpc.ClientConn
)

// ID maps the external ID used in Dgraph to the UID
type ID struct {
	Xid string `json:"xid,omitempty"`
	UID string `json:"uid,omitempty"`
}

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

	connection = conn
	dc := api.NewDgraphClient(connection)
	client = dgo.NewDgraphClient(dc)

	return nil
}

// Close terminates the Dgraph connection
func Close() {
	err := connection.Close()
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
		startTime: dateTime @index(hour) .
		endTime: dateTime @index(hour) .
		isService: bool .
		isPod: bool .
		isContainer: bool .
		isProc: bool .
		namespace: uid @reverse .
	`
	ctx := context.Background()
	err := client.Alter(ctx, op)

	return err
}

// GetUID returns the UID of the node in the Dgraph
// returns empty string if error has occurred
func GetUID(id string, nodeType string) string {
	query := `query Me($id:string, $nodeType:string) {
		getUid(func: eq(xid, $id)) @filter(has(` + nodeType + `)) {
			uid
		}
	}`

	ctx := context.Background()
	variables := make(map[string]string)
	variables["$nodeType"] = nodeType
	variables["$id"] = id

	resp, err := client.NewReadOnlyTxn().QueryWithVars(ctx, query, variables)
	if err != nil {
		log.Printf("failed to fetch UID from Dgraph %v", err)
		return ""
	}
	return unmarshalDgraphResponse(resp, id)
}

// ExecuteQuery given a query and it fetches and writes result into interface
func ExecuteQuery(query string, root interface{}) error {
	ctx := context.Background()

	resp, err := client.NewTxn().Query(ctx, query)
	if err != nil {
		log.Error(err)
		return err
	}

	err = json.Unmarshal(resp.Json, root)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// MutateNode mutates a Dgraph transaction
func MutateNode(data interface{}, mutateType string) (*api.Assigned, error) {
	bytes := utils.JSONMarshal(data)
	if bytes == nil {
		return nil, fmt.Errorf("Unable to marshal data: %v", data)
	}

	mu := &api.Mutation{
		CommitNow: true,
	}
	switch mutateType {
	case DELETE:
		mu.DeleteJson = bytes
	default:
		mu.SetJson = bytes
	}

	ctx := context.Background()
	return client.NewTxn().Mutate(ctx, mu)
}

// unmarshalDgraphResponse returns empty string if error has occurred
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
