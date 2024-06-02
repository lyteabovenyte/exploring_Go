package main

import (
	"encoding/csv"
	"fmt"
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

func writeRecs(recs []record) error {
	file, err := os.OpenFile("write_with_csv_package.csv", os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	//sort recs, before writing them to the file
	sort.Slice(recs, func(i, j int) bool {
		if len(recs[i]) > 0 && len(recs[j]) > 0 {
			return recs[i].last() < recs[j].last()
		}
		return false
	})

	for _, rec := range recs {
		if err := w.Write(rec); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	recs := []record{
		record{"Amir", "Alaeifar"},
		record{"Person", "B"},
		record{"Happy", "Person"},
		record{},
	}

	if err := writeRecs(recs); err != nil {
		panic(err)
	}

	data, err := os.ReadFile("write_with_csv_package.csv")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", data)

}
