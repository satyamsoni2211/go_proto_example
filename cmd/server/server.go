package main

import (
	"context"
	"github.com/satyamsoni2211/go_proto_example/pb"
	"github.com/satyamsoni2211/go_proto_example/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func unaryInterceptorfunc(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("method Path: %v", info.FullMethod)
	return handler(ctx, req)
}

func streamInterceptorFunc(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	log.Printf("method path :- %v", info.FullMethod)
	return handler(srv, ss)
}

func main() {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptorfunc),
		grpc.StreamInterceptor(streamInterceptorFunc),
	)
	listener, err := net.Listen("tcp", "[::]:8080")
	if err != nil {
		log.Fatalf("Cannot open listener on port 8080 :- %v", err)
	}
	FibServer := service.NewFibServer()
	pb.RegisterFibonacciServiceServer(server, FibServer)
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("Cannot server :- %v", err)
	}
}
