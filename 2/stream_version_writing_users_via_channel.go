package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
)

var users = []User{
	{Name: "amir", ID: 0},
	{Name: "morteza", ID: 1},
	{Name: "arad", ID: 2},
}

type User struct {
	Name string
	ID   int
}

func (u User) String() string {
	return fmt.Sprintf("%s:%d", u.Name, u.ID)
}

// writeUser starts a goroutine that will write User records received on 'in' and
// write them to 'w'. Any error encountered will write an error to the returned channel.
// If 'w' also implements io.Closer, we will .Close() on it.
func writeUser(ctx context.Context, w io.Writer, in chan User) chan error {
	errCh := make(chan error, 1)

	go func() {
		defer func() {
			if closer, ok := w.(io.Closer); ok {
				if err := closer.Close(); err != nil {
					// Try to put an error in our errCh,
					// otherwise, if we can't, just ignore it
					select {
					case errCh <- err:
					default:
					}
				}
			}
			close(errCh)
		}()

		writeLines := false

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case user, ok := <-in:
				if !ok {
					return
				}
				// this puts a carriage return before the next entry unless it
				// is the first entry.
				if writeLines {
					if _, err := w.Write([]byte("\n")); err != nil {
						errCh <- err
						return
					}
				}
				if _, err := w.Write([]byte(user.String())); err != nil {
					errCh <- err
					return
				}
				writeLines = true

			}
		}

	}()
	return errCh
}

func main() {
	// a byte.BufferP{} is an io.Writer. using this to avoid making files on disk
	// buffer types implements i/o package interface.
	buff := &bytes.Buffer{}

	// send our User records to be written via in.
	in := make(chan User, 1)

	// strat our goroutine for writing or User records to buff.
	errCh := writeUser(context.Background(), buff, in)

	for _, u := range users {
		select {
		case err := <-errCh:
			fmt.Println("had error: ", err)
			return
		// send another user to be written.
		case in <- u:
		}
	}

	// let our goroutine started by writeUser() know we have finished.
	close(in)

	if err := <-errCh; err != nil {
		fmt.Println("had error: ", err)
		return
	}

	// print our output buffer that was written
	io.Copy(os.Stdout, buff)
}
