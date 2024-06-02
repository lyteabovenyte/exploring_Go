package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
)

type record [][]byte

var comma = []byte(`,`)

var fakeFile = &bytes.Buffer{}

var fakeContent = `
Amir, Alaeifar
Happy, Person
Another, Happy
`

func (r record) validate() error {
	if len(r) != 2 {
		return errors.New("data format is incorrect")
	}
	return nil
}

func (r record) first() []byte {
	return r[0]
}

func (r record) last() []byte {
	return r[1]
}

func readRecs() ([]record, error) {
	scanner := bufio.NewScanner(fakeFile)

	var records []record
	for scanner.Scan() {
		line := scanner.Bytes()
		// skip the empty files
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		lineNum := 0
		var rec record = bytes.Split(line, comma)
		if err := rec.validate(); err != nil {
			return nil, fmt.Errorf("entry at line %d was invalid: %w", lineNum, err)
		}
		records = append(records, rec)
		lineNum++
	}
	return records, scanner.Err()
}

func main() {
	fakeFile.WriteString(fakeContent)

	recs, err := readRecs()
	if err != nil {
		panic(err)
	}
	for _, recs := range recs {
		fmt.Printf("FirstName: %s, LastName: %s\n", recs.first(), recs.last())
	}
}
