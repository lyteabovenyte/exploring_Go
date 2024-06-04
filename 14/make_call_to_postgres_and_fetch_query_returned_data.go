package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// a table in pg database.
// making a struct to fetch the query result in the variable using Scan.
type Category struct {
	categoryid   int
	categoryname string
	description  string
}

func getCategory(ctx context.Context, conn *sql.DB, categoryid int) (Category, error) {
	// const query = `SELECT "categoryname", "description" FROM categories WHERE "categoryid" = $1`

	c := Category{categoryid: categoryid}
	err := conn.QueryRowContext(
		ctx,
		"SELECT categoryname, description FROM categories WHERE categoryid = $1",
		categoryid,
	).Scan(&c.categoryname, &c.description)
	return c, err
}

func main() {
	dbURL := "postgres://username:password@localhost:5432/shop"

	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		fmt.Printf("connect to db error: %s\n", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cat, err := getCategory(ctx, conn, 1)
	if err != nil {
		fmt.Printf("error calling getCategory func: %v", err)
	} else {
		fmt.Printf("categoryname: %s\n", cat.categoryname)
		fmt.Printf("description: %s\n", cat.description)
	}
	cancel()
}
