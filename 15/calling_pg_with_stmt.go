package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Storage struct {
	conn            *sql.DB
	getCategoryStmt *sql.Stmt
}

type Category struct {
	categoryid   int
	categoryname string
	description  string
}

func NewStorage(ctx context.Context, conn *sql.DB) (*Storage, error) {

	stmt, err := conn.PrepareContext(ctx, "SELECT categoryname, description FROM categories WHERE categoryid = $1")
	if err != nil {
		return &Storage{}, fmt.Errorf("error createing stmt: %v", err)
	} else {
		return &Storage{
			conn:            conn,
			getCategoryStmt: stmt,
		}, nil
	}
}

func (s *Storage) getCategory(ctx context.Context, id int) (Category, error) {
	c := Category{categoryid: id}
	err := s.getCategoryStmt.QueryRow(id).Scan(&c.categoryname, &c.description)
	return c, err
}

func main() {
	dbURL := "postgres://username:password@localhost:5432/shop"

	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		fmt.Printf("error connecting to db: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	storage, err := NewStorage(ctx, conn)
	if err != nil {
		fmt.Printf("error on NewStorage function call: %v", err)
	}

	c, err := storage.getCategory(ctx, 1)
	if err != nil {
		fmt.Printf("error calling getCategory method on storage: %v", err)
	} else {
		fmt.Println("Categoy categoryname: ", c.categoryname)
		fmt.Println("Category description: ", c.description)
	}
	cancel()
}
