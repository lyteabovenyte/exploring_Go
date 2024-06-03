package main

import (
	"encoding/json"
	"io"
	"strings"
)

type Message struct {
	Name string
	Text string
}

// takes in a jsonStream and put the jsons into msgs and errors in errs channel
func decodeJson(jsonStram string) (chan Message, chan error) {
	msgs := make(chan Message, 1)
	errs := make(chan error, 1)

	reader := strings.NewReader(jsonStram)
	dec := json.NewDecoder(reader)
	// stream the jsonStream into msgs and errs channels
	go func() {
		defer close(msgs)
		defer close(errs)

		for {
			var m Message
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				errs <- err
				return
			}
			msgs <- m
		}
	}()
	return msgs, errs
}

// takes in and input message channel and wrtie JSON message into
// io.Writer output message to output to a stream
func encodeJson(in chan Message, output io.Writer) chan error {
	errs := make(chan error, 1)
	go func() {
		defer close(errs)
		enc := json.NewEncoder(output)
		for msg := range in {
			if err := enc.Encode(msg); err != nil {
				errs <- err
				return
			}
		}
	}()
	return errs
}
