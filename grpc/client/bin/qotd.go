package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/lyteabovenyte/exploring_go/grpc/client"
)

var (
	addr   = flag.String("addr", "127.0.0.1:2562", "The address of the server.")
	author = flag.String("author", "", "The author whose quote to get")
)

func main() {
	flag.Parse()

	c, err := client.New(*addr)
	if err != nil {
		panic(err)
	}

	a, q, err := c.QOTD(context.Background(), *author)
	if err != nil {
		panic(err)
	}

	fmt.Println("Author: ", a)
	fmt.Printf("Quote of the Day: %q\n", q)
}
