package main

import (
	"ylog/core"
	"strconv"
	"sync"
)

func main() {
	var w sync.WaitGroup
	for i := 0; i < 10000; i++ {
		i:=i
		w.Add(1)
		go func() {
			for j:=0;j<100;j++  {
				ylog.Info("xixi"+strconv.Itoa(i), "./sql/log.txt")
			}
			w.Done()
		}()
	}
	w.Wait()
}
