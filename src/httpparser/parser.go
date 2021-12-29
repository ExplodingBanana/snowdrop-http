package httpparser

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var errRequestSyntax = errors.New("request syntax error")
var errNoSplitter = errors.New("no splitter was found")

type IProtocol interface {
	OnMessageBegin()
	OnMethod([]byte)
	OnPath([]byte)
	OnProtocol([]byte)
	OnHeadersBegin()
	OnHeader(string, string)
	OnHeadersComplete()
	OnBody([]byte)
	OnMessageComplete()
}

type HTTPParser struct {
	_               struct{}
	CurrentState    ParsingState
	Protocol        IProtocol
	
	сurrentSplitter []byte
	buffer        []byte
	contentLength uint
}

func (parser *HTTPParser) Feed(data []byte) (completed bool, err error) {
	if parser.сurrentSplitter == nil {
		return false, errNoSplitter
	}

	for _, v := range data {
		if v == 10 {
			parser.CurrentState +=1
		}
	}

	parser.buffer = append(parser.buffer, data...)
	fmt.Println(parser.buffer)

	switch parser.CurrentState {
	case MessageCompleted:
		parser.CurrentState = Ready
		return true, nil
	}

	return true, nil
}

func SplitBytes(src, splitBy []byte) [][]byte {
	if len(src) == 0 {
		return [][]byte{}
	}

	var splited [][]byte
	var afterPrevSplitBy uint
	var skipIters int
	lookForward := len(splitBy)

	for index := range src[:len(src)-lookForward] {
		if skipIters > 0 {
			skipIters--
			continue
		}

		if bytes.Equal(src[index:index+lookForward], splitBy) {
			splited = append(splited, src[afterPrevSplitBy:index])
			afterPrevSplitBy = uint(index + lookForward)
			skipIters = lookForward
		}
	}

	if len(splited) == 0 {
		splited = append(splited, src)
	} else if bytes.HasSuffix(src, splitBy) {
		// if source ends with splitter, we must add pending
		// shit without counting splitter in the end
		splited = append(splited, src[afterPrevSplitBy:len(src)-lookForward])
	} else {
		// or add pending shit, but with counting everything in the end
		splited = append(splited, src[afterPrevSplitBy:])
	}

	return splited
}

func parseHeaders(rawHeaders []byte) (parsedHeaders map[string]string, err error) {
	headers := map[string]string{}

	for _, rawHeader := range SplitBytes(rawHeaders, []byte(CLRF)) {
		key, value, err := parseHeader(rawHeader)

		if err != nil {
			return nil, err
		}

		headers[*key] = *value
	}

	return headers, nil
}

func parseHeader(headersBytesString []byte) (key *string, value *string, err error) {
	for index, char := range headersBytesString {
		if char == ':' {
			key := string(headersBytesString[:index])
			value := string(headersBytesString[index+1:])

			value = strings.TrimPrefix(value, " ")

			return &key, &value, nil
		}
	}

	return nil, nil, errNoSplitter
}
