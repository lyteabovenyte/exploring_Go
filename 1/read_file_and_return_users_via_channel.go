package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const filePath = "/Users/lyteatnyte/Developer/Go/Go-for-DevOps-test/read_file_and_return_users_via_channel/users.txt"

// define an init function to populate the users on startup.
func init() {
	content := []byte("amir:0\nmorteza:1\narad:3")

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		panic(err)
	}
}

// define a User struct type
type User struct {
	Name string
	ID   int

	err error
}

// implementing the Stringer method for out User type.
func (u User) String() string {
	return fmt.Sprintf("%s:%d", u.Name, u.ID)
}

// getUser takes a string that sould be formatted to [user]:[id]
func getUser(s string) (User, error) {
	sp := strings.Split(s, ":")
	if len(sp) != 2 {
		return User{}, fmt.Errorf("record(%s) was not in the correct format", s)
	}
	id, err := strconv.Atoi(sp[1])
	if err != nil {
		return User{}, fmt.Errorf("record(%s) has non-numeric ID", s)
	}
	return User{Name: strings.TrimSpace(sp[0]), ID: id}, nil
}

func decodeUser(ctx context.Context, r io.Reader) chan User {
	ch := make(chan User, 1)

	go func() {
		defer close(ch)

		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if ctx.Err() != nil {
				ch <- User{err: ctx.Err()}
				return
			}
			u, err := getUser(scanner.Text())
			if err != nil {
				u.err = err
				ch <- u
				return
			}
			// unless, everything was fine
			ch <- u
		}
	}()
	// return the channel
	return ch
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	//start decoding the file one line at a time
	ch := decodeUser(ctx, f)

	for u := range ch {
		if u.err != nil {
			panic(err)
		}
		fmt.Println(u)
	}
}
