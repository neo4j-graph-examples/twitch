// go mod init main
// go run example.go
package main

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"reflect"
)

func main() {
	results, err := runQuery("neo4j://<HOST>:<BOLTPORT>", "neo4j", "<USERNAME>", "<PASSWORD>")
	if err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}

func runQuery(uri, database, username, password string) (_ []string, err error) {
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	defer func() { err = handleClose(ctx, driver, err) }()
	query := "	MATCH (s:Stream{name:$streamer})<-[:VIP|:MODERATOR]-(user)
	RETURN user.name as moderator_or_vip LIMIT 20
	params := map[string]any{"streamer": "ludwig"}
	result, err := neo4j.ExecuteQuery(ctx, driver, query, params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(database),
		neo4j.ExecuteQueryWithReadersRouting())
	if err != nil {
		return nil, err
	}
	moderator_or_vips := make([]string, len(result.Records))
	for i, record := range result.Records {
		// this assumes all actors have names, hence ignoring the 2nd returned value
		name, _, err := neo4j.GetRecordValue[string](record, "moderator_or_vip")
		if err != nil {
			return nil, err
		}
		moderator_or_vips[i] = name
	}
	return moderator_or_vips, nil
}

func handleClose(ctx context.Context, closer interface{ Close(context.Context) error }, previousError error) error {
	err := closer.Close(ctx)
	if err == nil {
		return previousError
	}
	if previousError == nil {
		return err
	}
	return fmt.Errorf("%v closure error occurred:\n%s\ninitial error was:\n%w", reflect.TypeOf(closer), err.Error(), previousError)
}
