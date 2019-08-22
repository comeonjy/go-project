package main

import (
	"fmt"
	"time"
)

func main()  {
	ch:=make(chan func())
	go func() {
		fmt.Println(<-ch)
	}()
	go func() {
		close(ch)
	}()
	time.Sleep(1e9)
}