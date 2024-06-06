package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
)

// defining a custom flag
type URLValue struct {
	URL *url.URL
}

// by defining a custom flag, we should implement flag.Value methods
// to satisfy it's interface
// there are two methods on type Value interface for the flag
// 1. String()  2.Set(s string)
func (v URLValue) String() string {
	if v.URL != nil {
		return v.URL.String()
	}
	return ""
}

// by defining a Set() method on a type, we can read in any custom values
func (v URLValue) Set(s string) error {
	if u, err := url.Parse(s); err != nil {
		return err
	} else {
		*v.URL = *u
	}
	return nil
}

var u = &url.URL{}

// basic flag error handling
var (
	useProd = flag.Bool("prod", true, "Use a production endpoint")
	useDev  = flag.Bool("dev", false, "Use a development endpoint")
	help    = new(bool)
)

func init() {
	flag.Var(&URLValue{u}, "url", "URL to parse")
	flag.BoolVar(help, "help", false, "Display help text")
	flag.BoolVar(help, "h", false, "Display help text")
}

func main() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	switch {
	case *useDev && *useProd:
		log.Println("Error: --prod and --dev cannot both be set")
		flag.PrintDefaults()
		os.Exit(1)
	case !(*useProd && *useDev):
		log.Println("Error: either --prod or --dev must be set")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if reflect.ValueOf(*u).IsZero() {
		panic("did you pass the url?")
	}
	fmt.Printf("{scheme: %q, host: %q, path: %q}", u.Scheme, u.Host, u.Path)
}
