package main

import (
	"context"
	"fmt"

	"github.com/everywan/identifier"
	"github.com/everywan/identifier/pb"
)

func main() {
	general()
}

func general() {
	opts := &identifier.SnowflakeClientOps{
		Address: "127.0.0.1:10001",
	}
	client, err := identifier.NewSnowflakeClient(opts)
	if err != nil {
		fmt.Println("can't connect snowflake server")
		return
	}
	req := &pb.Request{}
	fmt.Println(client.Generate(context.Background(), req))
}
