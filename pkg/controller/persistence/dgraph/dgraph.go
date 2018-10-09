package dgraph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

var Client *dgo.Dgraph
var Connection *grpc.ClientConn

func init() {
	Open("127.0.0.1:9080")
	err := CreateSchema()
	if err != nil {
		fmt.Println("Error while creating schema ", err)
	}
}

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

func Close() {
	Connection.Close()
}

func CreateSchema() error {
	op := &api.Operation{}
	op.Schema = `
		name: string @index(term) .
		xid:  string @index(term) .
		isService: string @index(term) .
		isPod: string @index(term) .
	`
	ctx := context.Background()
	err := Client.Alter(ctx, op)

	return err
}

func GetUId(dg *dgo.Dgraph, id string, nodeType string) (string, error) {

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
		return "", err
	}

	type Root struct {
		IDs []ID `json:"getUid"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return "", err
	}

	if len(r.IDs) == 0 {
		return "", fmt.Errorf("id %s is not in dgraph", id)
	}

	return r.IDs[0].UID, nil
}

func MutateNode(dg *dgo.Dgraph, n []byte) error {
	mu := &api.Mutation{
		CommitNow: true,
	}

	mu.SetJson = n
	ctx := context.Background()
	_, err := dg.NewTxn().Mutate(ctx, mu)

	return err
}
