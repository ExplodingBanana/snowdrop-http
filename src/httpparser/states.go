package httpparser

type ParsingState uint8

const (
	Ready ParsingState = iota + 1
	Method
	Path
	Protocol
	Headers
	Body
	MessageCompleted
)
