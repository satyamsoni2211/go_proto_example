package main

import (
	"context"
	"flag"
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
	address := flag.String("address", "[::]:8080", "port to start server")
	flag.Parse()
	server := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptorfunc),
		grpc.StreamInterceptor(streamInterceptorFunc),
	)
	listener, err := net.Listen("tcp", *address)
	log.Printf("Starting server on %v", *address)
	if err != nil {
		log.Fatalf("Cannot open listener on port 8080 :- %v", err)
	}
	FibServer := service.NewFibServer()
	ImageServer := service.NewImageServer("img", make(map[string]string))

	//registering Fibonacci server
	pb.RegisterFibonacciServiceServer(server, FibServer)
	//registering image server
	pb.RegisterUploadImageServiceServer(server, ImageServer)

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("Cannot server :- %v", err)
	}
}
