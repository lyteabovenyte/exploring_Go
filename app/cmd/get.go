/*
Copyright © 2024 Amir Alaeifar lyteabovenyte@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lyteabovenyte/exploring_go/grpc/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		const devAddr = "127.0.0.1:3450"

		fs := cmd.Flags()

		addr := mustString(fs, "addr") // using --addr flag's value.

		if mustBool(fs, "dev") {
			addr = devAddr
		}

		c, err := client.New(addr)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}
		a, q, err := c.QOTD(cmd.Context(), mustString(fs, "author")) // using an --author flag value which defaults to an empty string
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}

		switch {
		case mustBool(fs, "json"): // using a --json flags to determine whether the output should be in json.
			b, err := json.Marshal(
				struct {
					Author string
					Quote  string
				}{a, q},
			)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", b)
		default:
			fmt.Println("Author: ", a)
			fmt.Println("Quote: ", q)
		}
	},
}

func mustString(fs *pflag.FlagSet, name string) string {
	v, err := fs.GetString(name)
	if err != nil {
		panic(err)
	}
	return v
}

func mustBool(fs *pflag.FlagSet, name string) bool {
	v, err := fs.GetBool(name)
	if err != nil {
		panic(err)
	}
	return v
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolP("dev", "d", false,
						 "Users the dev server instead of prod")
	getCmd.Flags().String("addr", "127.0.0.1:2562",
						"Set the QOTD server to use, defaults to production")
	getCmd.Flags().StringP("author", "a", "",
							"Specify the author to get a quote for")
	getCmd.Flags().Bool("json", false,
						 "Output is in JSON format")	

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
