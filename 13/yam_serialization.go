package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Jobs []Job
}

type Job struct {
	Name     string
	Interval time.Duration
	Cmd      string
}

func main() {
	c := Config{
		Jobs: []Job{
			{
				Name:     "clear tmp",
				Interval: 24 * time.Hour,
				Cmd:      "rm -rf " + os.TempDir(),
			},
		},
	}

	b, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}
