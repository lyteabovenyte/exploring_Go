// if we want to read and sort the records
// and also convert them back to csv in another file.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
)

type record []string

func (r record) first() string {
	return r[0]
}

func (r record) last() string {
	return r[1]
}

func (r record) csv() []byte {
	b := bytes.Buffer{}
	for _, field := range r {
		b.WriteString(field + ",")
	}
	b.WriteString("\n")
	return b.Bytes()
}

func writeRecs(recs []record) error {
	file, err := os.OpenFile("data-sorted.csv", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// sort the recs by their last name
	sort.Slice(recs, func(i, j int) bool { return recs[i].last() < recs[j].last() })

	for _, rec := range recs {
		_, err := file.Write(rec.csv())
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func main() {
	recs := []record{
		record{"Amir", "Alaeifar"},
		record{"Happy", "Person"},
		record{"Another", "One"},
	}

	if err := writeRecs(recs); err != nil {
		panic(err)
	}

	data, err := os.ReadFile("data-sorted.csv")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", data)
}
