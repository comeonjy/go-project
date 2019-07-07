package engine

import (
	"ccmouse/fetcher"
	"log"
)

type ConcurrentEngine struct {
	Scheduler Scheduler
	WorkerCount int
}

type Scheduler interface {
	Submit(Request)
	ConfigureMasterWorkerChan(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request)  {
	in :=make(chan Request)
	out :=make(chan ParseResult)
	e.Scheduler.ConfigureMasterWorkerChan(in)

	for i:=0;i<e.WorkerCount;i++{
		createWorker(in,out)
	}

	for _,r:=range seeds{
		e.Scheduler.Submit(r)
	}

	for{
		result:=<-out
		for _,item:=range result.Items{
			log.Printf("Got item: %v",item)
		}
		for _,request:=range result.Requests{
			e.Scheduler.Submit(request)
		}
	}
}

func createWorker(in chan Request,out chan ParseResult)  {
	go func() {
		for  {
			request:=<-in
			result,err:=worker(request)
			if err != nil {
				continue
			}
			out<-result
		}
	}()
}

func worker(r Request) (ParseResult,error) {
	body,err:=fetcher.Fetcher(r.Url)
	if err!=nil {
		log.Printf("Resp Error:%s",err)
		return ParseResult{},err
	}
	return r.ParserFunc(body),nil
}
