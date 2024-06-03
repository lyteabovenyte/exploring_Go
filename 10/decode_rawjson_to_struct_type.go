// decode raw json with brackets in the begining and trainling bracket
// and carriage between each message
// using dec.Token() to remove them.
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const jsonStream = `[
	{"Name": "amir", "Text": "Hello Happy man"},
	{"Name": "gennday", "Text": "Hello senior"}
]`

type Message struct {
	Name string
	Text string
}

func main() {
	reader := strings.NewReader(jsonStream) // converting to io.Reader interface

	dec := json.NewDecoder(reader)

	_, err := dec.Token() // reads first [
	if err != nil {
		panic(fmt.Errorf("outer [ is missing"))
	}

	for dec.More() {
		var m Message
		//decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", m)
	}

	_, err = dec.Token() // read last ]
	if err != nil {
		panic(fmt.Errorf("final ] is missing"))
	}
}
