package parser

import (
	"ccmouse/engine"
	"regexp"
)

func ParseCity(contents []byte) engine.ParseResult {
	//compiles := regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[a-zA-Z0-9]+)"[^>]*>([^<].)</a>`)
	//submatch := compiles.FindAllSubmatch(contents, -1)
	var result engine.ParseResult
	//for _, m := range submatch {
	//	result.Items=append(result.Items,string(m[2]))
	//	result.Requests=append(result.Requests,engine.Request{
	//		Url:        string(m[1]),
	//		ParserFunc: ParseCity,
	//	})
	//}


	member_rg := regexp.MustCompile(`<a href="(http://[a-zA-Z0-9]+.zhenai.com/u/[0-9]+)"[^>]*>([^<].)</a>`)
	member_submatch := member_rg.FindAllSubmatch(contents, -1)
	for _, m := range member_submatch {
		result.Items=append(result.Items,string(m[2]))
		result.Requests=append(result.Requests,engine.Request{
			Url:        string(m[1]),
			ParserFunc: ParseProfile,
		})
	}

	return result
}
