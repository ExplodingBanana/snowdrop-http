package main

import (
	"fmt"
	"os"

	"gtihub.com/floordiv/snowdrop/src/httpparser"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var parser = httpparser.HTTPParser{CurrentState: httpparser.Ready}

	data, err := os.ReadFile("temp")
	check(err)
	parser.Feed(data)
	fmt.Printf("parser: %v\n", parser)
}
