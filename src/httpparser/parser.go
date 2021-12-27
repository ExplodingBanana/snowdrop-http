package httpparser

import (
	"bytes"
	"errors"
	"strings"
)

var errRequestSyntax = errors.New("request syntax error")
var errNoSplitter = errors.New("no splitter was found")

type IProtocol interface {
	OnReady()
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
	_            struct{}
	CurrentState ParsingState
	Protocol     IProtocol

	contentLength   uint
	currentSplitter []byte
}

func (parser *HTTPParser) Feed(data []byte) (completed bool, err error) {
	if parser.CurrentState == MessageCompleted {
		parser.CurrentState = Ready
		parser.Protocol.OnReady()
		return true, nil
	}
	parser.contentLength = uint(len(data))
	parser.currentSplitter = []byte(CLRF)

	headers := map[string]string{}

	parser.Protocol.OnMessageBegin()

	var splitted_data = bytes.Split(data, parser.currentSplitter)
	for index, line := range splitted_data {
		if index == 0 {
			var splitted_line = bytes.Split(line, []byte(" "))
			parser.Protocol.OnMethod(splitted_line[0])
			parser.Protocol.OnPath(splitted_line[1])
			parser.Protocol.OnMethod(splitted_line[2])
			continue
		}
		if parser.CurrentState == Headers {
			if bytes.Equal(line, parser.currentSplitter) {
				parser.CurrentState = Body
				continue
			}
			key, value, _ := parseHeader(line)
			parser.Protocol.OnHeader(*key, *value)
			headers[*key] = *value
			continue
		}
		if bytes.Equal(line, parser.currentSplitter) {
			return false, errRequestSyntax
		}
		if bytes.Equal(line, parser.currentSplitter) {
			parser.CurrentState = Body
		}
	}

	/*
		for index, char := range data {
			if char == parser.currentSplitter {
				if data[index-1] == parser.currentSplitter {
					if parser.CurrentState != Headers {
						return false, errRequestSyntax
					}

					parser.CurrentState = Body
					parser.Protocol.OnHeadersComplete()

				}
			}
		}
	*/

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

/*
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
*/

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
