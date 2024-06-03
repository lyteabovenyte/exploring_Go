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

func YamlMarshal(jobconf []Job) (b []byte, err error) {
	c := Config{
		Jobs: jobconf,
	}

	M, err := yaml.Marshal(c)
	if err != nil {
		return []byte{}, fmt.Errorf("an error happend: %v", err)
	}
	return M, nil
}

func YamlUnmarshal(data []byte) (c Config, err error) {
	cfg := Config{}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("err: %v", err)
	}
	return cfg, nil
}

func main() {
	var examplejob = []Job{
		{
			Name:     "clear tmp",
			Interval: 24 * time.Hour,
			Cmd:      "rm -rf " + os.TempDir(),
		},
	}

	var exampledata = []byte(`
	jobs:
		- name: Clear tmp
		  interval: 24h0m0s
		  whatever: is not in the Job type
		  cmd: rm -rf /tmp
	`)
	fmt.Println("________________")
	fmt.Println("MARSHALING:\n")

	b, err := YamlMarshal((examplejob))
	if err != nil {
		fmt.Println("error on marshaling: ", err)
	} else {
		fmt.Printf("%s\n", b)
	}

	fmt.Println("________________")
	fmt.Println("UNMARSHALING:\n")

	c, err := YamlUnmarshal(exampledata)
	if err != nil {
		fmt.Println("error on Unmarshalling: ", err)
	} else {
		for _, job := range c.Jobs {
			fmt.Println("Name: ", job.Name)
			fmt.Println("Interval: ", job.Interval)
		}
	}
}
