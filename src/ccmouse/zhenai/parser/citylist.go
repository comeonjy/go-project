package parser

import (
	"ccmouse/engine"
	"regexp"
)

func ParseCityList(contents []byte) engine.ParseResult {
	compile := regexp.MustCompile(`<a href="(http://www.zhenai\.com/zhenghun/[a-zA-z0-9]+)" [^>]*>([^<]+)</a>`)
	submatchs := compile.FindAllSubmatch(contents, -1)
	result:=engine.ParseResult{}
	limit:=1
	for _,m:=range submatchs {
		result.Items = append(result.Items,string(m[2]))
		result.Requests = append(result.Requests,engine.Request{
			Url:        string(m[1]),
			ParserFunc: ParseCity,
		})
		limit--
		if limit<=0 {
			break
		}
	}
	return result
}
