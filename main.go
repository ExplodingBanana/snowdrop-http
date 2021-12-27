package main

import (
	"fmt"
	"os"

	"gtihub.com/floordiv/snowdrop/src/httpparser"
)

func check(e error) {
	if e != nil {
		print("fok")
		panic(e)
	}
}

func main() {
	var parser = httpparser.HTTPParser{CurrentState: httpparser.Ready}

	data, err := os.ReadFile("test")
	check(err)
	fmt.Println(data)
	fmt.Printf("parser: %v\n", parser)
	completed, err := parser.Feed(data)
	if completed {
		return
	}
	check(err)
}
