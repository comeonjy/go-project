package main

import (
	"ccmouse/engine"
	"ccmouse/scheduler"
	"ccmouse/zhenai/parser"
)

func main()  {
	e:=engine.ConcurrentEngine{
		Scheduler:&scheduler.SimpleScheduler{},
		WorkerCount:10,
	}
	e.Run(engine.Request{
		Url:        "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}


