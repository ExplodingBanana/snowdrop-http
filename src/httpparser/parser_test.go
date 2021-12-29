package httpparser

import (
	"fmt"
	"testing"
)

func TestHTTPParser_Feed(t *testing.T) {
	var parser = HTTPParser{
		CurrentState:    0,
		Protocol:        nil,
		сurrentSplitter: []byte("\n"),
		buffer:          []byte{},
		contentLength:   0,
	}
	fmt.Println("Current case: // regular text, no splitters")
	parser.Feed([]byte("GET / HTTP"))
	fmt.Println("Current case: // trailing spliiter")
	parser.Feed([]byte("/1.1\r\n"))
	fmt.Println("Current case: //splitter between text")
	parser.Feed([]byte("/1.1\r\nHost"))
	fmt.Println("Current case: // 2 trailing splitters")
	parser.Feed([]byte("dv.mz.org\r\n\r\n"))
	fmt.Println("Current case: // headers2body")
	parser.Feed([]byte("dv.mz.org\r\n\r\nbody"))
	fmt.Println("Current case: // leading splitter")
	parser.Feed([]byte("\r\nHost: "))
}
