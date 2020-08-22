package service

import (
	"github.com/satyamsoni2211/go_proto_example/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type FibonacciServer interface {
	CalculateFibonacci(*pb.FibonacciRequest, pb.FibonacciService_CalculateFibonacciServer) error
}

type FibServer struct {
}

func (f FibServer) CalculateFibonacci(
	request *pb.FibonacciRequest,
	stream pb.FibonacciService_CalculateFibonacciServer) error {
	number := request.GetNumber()
	for i := range calcFib(number) {
		res := &pb.FibonacciResponse{
			Result: i,
		}
		err := stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Internal, "Cannot generate fibonacci: %v", err)
		}
	}
	return nil
}

func NewFibServer() *FibServer {
	return &FibServer{}
}

func calcFib(n uint32) <-chan uint32 {
	out := make(chan uint32)

	go func() {
		x, y := 0, 1
		counter := 0
		defer close(out)
		for {
			select {
			case out <- uint32(x):
				x, y = y, x+y
				counter += 1
				log.Println("Incremented counter to ", counter)
				if uint32(counter) > n {
					return
				}
			}
		}
	}()
	return out
}
