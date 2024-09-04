package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	num := "۰۹۳۷۲۸۲۸۰۰۸"

	mapper := map[rune]rune{
		'۱': '1',
		'۲': '2',
		'۳': '3',
		'۴': '4',
		'۵': '5',
		'۶': '6',
		'۷': '7',
		'۸': '8',
		'۹': '9',
		'۰': '0',
	}

	converted_string := strings.Map(func(r rune) rune {
		if c, ok := mapper[r]; ok {
			return c
		}
		return r
	}, num)

	converted_number, err := strconv.Atoi(converted_string)
	if err != nil {
		fmt.Println("error converting", converted_string)
	}
	fmt.Println(converted_number)
}
