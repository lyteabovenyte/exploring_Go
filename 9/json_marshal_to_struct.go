package main

import (
	"encoding/json"
	"fmt"
)

type Record struct {
	Name string `json:"user_name"`
	User string `json:"user"`
	ID   int
	Age  int `json:"-"` // we don't want this field to  be marshaled.
}

func main() {
	rec := Record{
		Name: "Amir Alaeifar", // this field would be converted to user_name as we mentioned in the field tag.
		User: "lyteabovenyte", // converted to user field as mentioned in the field tag.
		ID:   23,
		Age:  24,
	}

	b, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}
