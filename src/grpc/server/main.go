package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	. "pb"
)

type Server struct {
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Errorf("tcp监听错误：%s", err)
	}
	server := grpc.NewServer()
	RegisterGreeterServer(server, &Server{})
	if err := server.Serve(lis); err != nil {
		fmt.Errorf("%s", err)
	}
}

func (s *Server) SayHello(c context.Context, req *HelloRequest) (*HelloReply, error) {
	fmt.Println(req)
	return &HelloReply{
		Message: req.Name,
	}, nil
}
