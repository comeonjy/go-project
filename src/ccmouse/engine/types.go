package engine

type Request struct {
	Url string
	ParserFunc func([]byte) ParseResult
}

type ParseResult struct {
	Requests []Request
	Items []interface{}
}

func NilPaeser(b []byte) ParseResult{
	return ParseResult{}
}