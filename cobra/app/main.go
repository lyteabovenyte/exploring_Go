package main

import (
	"context"
	
	"github.com/lyteabovenyte/exploring_go/cobra/app/cmd"
)

func main() {
	ctx := context.Background()

	cmd.Execute(ctx)
}