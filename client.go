package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"github.com/satyamsoni2211/go_proto_example/pb"
)

func main() {
	con, err := grpc.Dial("0.0.0.0:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to server %v", err)
	}
	client := pb.NewFibonacciServiceClient(con)
	req := &pb.FibonacciRequest{
		Number: 10,
	}
	stream, err := client.CalculateFibonacci(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not receive response from %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("End of stream")
			return
		}
		if err != io.EOF && err != nil {
			log.Fatalf("Unable to receive stream response from server %v", err)
		}
		fmt.Println(res.GetResult())
	}
}
