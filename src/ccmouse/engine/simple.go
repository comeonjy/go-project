package engine

import (
	"ccmouse/fetcher"
	"fmt"
	"log"
)

type SimpleEngine struct {

}

func (e SimpleEngine)Run(seeds ...Request)  {
	var requests []Request
	for _,r:=range seeds  {
		requests=append(requests,r)
	}
	for len(requests)>0 {
		r:=requests[0]
		requests = requests[1:]
		fmt.Printf("Fetching:%s \n",r.Url)
		parseResult,err:=e.worker(r)
		if err!=nil {
			log.Println(err)
			continue
		}
		requests = append(requests,parseResult.Requests...)
		for _,item:=range parseResult.Items {
			fmt.Printf("Got item %v \n",item)
		}
	}
}

func (SimpleEngine) worker(r Request) (ParseResult,error) {
	body,err:=fetcher.Fetcher(r.Url)
	if err!=nil {
		log.Print("Resp Error:",err)
		return ParseResult{},err
	}
	return r.ParserFunc(body),nil
}
