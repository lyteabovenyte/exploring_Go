// conversion after reading the whole file in csv files data format
// When doing basic CSV manipulation, sometimes,
// it is easier to simply split data using a carriage return and 
// then split the line based on a comma or other separator.
// Let's say we have a CSV file
// representing first and last names and break that CSV file into records:
package main

import (
	"fmt"
	"strings"
)

// represents a file's content.
const fakeContent = `
Amir, Alaeifar
Person, B
Another, Person
`

type record []string

// implementing methods for the record type
// validating if the csv line had the correct number of enteries
func (r record) validate() error {
	if len(r) != 2 {
		return fmt.Errorf("data format is incorrect")
	}
	return nil
}

// return first element in record
func (r record) first() string {
	return r[0]
}

// return last element in record
func (r record) last() string {
	return r[1]
}

// readRecs reads a file in csv format representing the records and returns the records.
func readRecs() ([]record, error) {

	lines := strings.Split(fakeContent, "\n") // split lines

	var records []record
	for i, line := range lines {
		// skip the empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}
		var rec record = strings.Split(line, ",")
		if err := rec.validate(); err != nil {
			return nil, fmt.Errorf("entry at line %d was invalid: %w", i, err)
		}
		records = append(records, rec)
	}
	return records, nil
}


func main() {
	recs, err := readRecs()
	if err != nil {
		panic(err)
	}
	for _, rec := range recs {
		fmt.Printf("Name: %s, Family: %s\n", rec.first(), rec.last())
	}
}
