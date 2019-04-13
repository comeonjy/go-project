package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main(){
	listerner,err:=net.Listen("tcp","localhost:8000")
	if err!=nil{
		log.Print(err)
	}
	for {
		conn,err:=listerner.Accept()
		if err!=nil {
			log.Print(err)
		}
		fmt.Println("监听ing")
		go handerConn(conn)
	}
}

func handerConn(c net.Conn)  {
	defer c.Close()
	for {
		_,err:=io.WriteString(c,time.Now().Format("15:04:05\n"))
		if err!=nil {
			fmt.Println("处理完成")
			return
		}
		time.Sleep(1*time.Second)
	}
}