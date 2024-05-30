package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type record []string

func (r record) validate() error {
	if len(r) != 2 {
		return fmt.Errorf("data format is incorrect.")
	}
	return nil
}

func (r record) first() string {
	return r[0]
}

func (r record) last() string {
	return r[1]
}

func readRecs() ([]record, error) {
	file, err := os.Open("data.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = 2

	var recs []record

	for {
		data, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		rec := record(data)
		recs = append(recs, rec)
	}
	return recs, nil
}

const fakeContent = `
Amir, Alaeifar
Happy, Person
Another, One
`

func main() {
	file, err := os.OpenFile("data.csv", os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	if _, err := file.Write([]byte(fakeContent)); err != nil {
		panic(err)
	}

	recs, err := readRecs()
	if err != nil {
		panic(err)
	}

	for _, rec := range recs {
		fmt.Printf("First: %s, Last: %s\n", rec.first(), rec.last())
	}
}
