package main

import (
	"bufio"
	"context"
	"flag"
	"github.com/satyamsoni2211/go_proto_example/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

//CalcFibonacci Function to calculate fibonacci using
//server stream
//this is example of server streaming
func CalcFibonacci(client pb.FibonacciServiceClient) {
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
		log.Println(res.GetResult())
	}
}

//UploadImage Function to upload image to server using streaming
//this is example of client streaming
func UploadImage(imageClient pb.UploadImageServiceClient) {
	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Details_{
			Details: &pb.UploadImageRequest_Details{
				Type: "jpeg",
				Name: "download.jpeg",
			},
		},
	}
	stream, err := imageClient.UploadImage(context.Background())
	if err != nil {
		log.Fatalf("Cannot fetch stream details from server: %v", err)
	}
	err = stream.Send(req)
	if err != nil {
		log.Fatalf("Cannot send image stream response to server %v", err)
	}
	file, err := os.Open("download.jpeg")
	if err != nil {
		log.Fatalf("Cannot open file : %v", err)
	}

	//creating buffer and reader for reading bytes
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1<<10)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			log.Println("Done reading data")
			break
		}
		if err != nil && err != io.EOF {
			log.Fatalf("Cannot read data into byte %v", err)
		}
		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_Chunk{
				Chunk: buffer,
			},
		}
		err = stream.Send(req)

		if err != nil {
			log.Fatalf("Cannot send image data to server : %v", err)
		}
		log.Printf("read %d bytes into buffer \n", n)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Cannot receive final response from client %v", err)
	}

	log.Println(res.GetId(), res.GetSize())
}

func main() {

	//getting adress from the user for the server
	address := flag.String("address", "[::]:8080", "port to connect to client")
	flag.Parse()

	log.Printf("Connecting client at %v", *address)

	//creating connection for the grpc
	con, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to server %v", err)
	}

	//client creating for Fibonacci
	client := pb.NewFibonacciServiceClient(con)
	CalcFibonacci(client)

	//client creation for image upload
	//imageClient := pb.NewUploadImageServiceClient(con)
	//UploadImage(imageClient)

}
