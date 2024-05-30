// stram the content using bufio package and bufio scanners.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// fakeFile represents a file's content.
var fakeFile = &bytes.Buffer{}

const fakeContent = `
Amir, Alaeifar
Happy, Person
motivated, Person
`

type record []string

func (r record) validate() error {
	if len(r) != 2 {
		return fmt.Errorf("data format is incorrect")
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

	scanner := bufio.NewScanner(fakeFile)

	var records []record
	lineNum := 0
	for scanner.Scan() { // Scan's default is scanning by line.
		line := scanner.Text()
		// skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		var rec record = strings.Split(line, ",") // split by delimiter on the file which is ,
		if err := rec.validate(); err != nil {
			return nil, fmt.Errorf("entry at line %d was invalid: %w", lineNum, err)
		}
		records = append(records, rec)
		lineNum++
	}
	return records, scanner.Err()
}

func main() {
	// create our fakeFile
	fakeFile.WriteString(fakeContent)

	recs, err := readRecs()
	if err != nil {
		panic(err)
	}
	for _, rec := range recs {
		fmt.Printf("Name: %s, Last: %s\n", rec.first(), rec.last())
	}
}
