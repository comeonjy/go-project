package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	. "pb"
)

func main()  {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err!=nil {
		fmt.Errorf("%s",err)
	}
	defer conn.Close()
	client:=NewGreeterClient(conn)
	reply,err:=client.SayHello(context.Background(),&HelloRequest{
		Name: "i am client",
	})
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(reply.Message)

}
