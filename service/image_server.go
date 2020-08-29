package service

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/satyamsoni2211/go_proto_example/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"log"
)

type ImageServiceServer interface {
	UploadImage(pb.UploadImageService_UploadImageServer) error
}

type ImageServer struct {
	destination string
	store       map[string]string
}

func NewImageServer(destination string, store map[string]string) *ImageServer {
	return &ImageServer{destination: destination, store: store}
}

func (i *ImageServer) SaveImage(name string, data []byte) (string, error) {
	fileName := fmt.Sprintf("%v/%v", i.destination, name)
	err := ioutil.WriteFile(fileName, data, 0644)
	id := uuid.New().String()
	i.store[id] = name
	if err != nil {
		return "", err
	}
	return id, nil
}

func (i ImageServer) UploadImage(stream pb.UploadImageService_UploadImageServer) error {
	resp, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "Cannot recieve stream response from client: %v", err)
	}
	details := resp.GetDetails()
	data := new(bytes.Buffer)
	writer := bufio.NewWriter(data)
	size := 0
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return status.Errorf(codes.Unknown, "Cannot receive chunks from stream %v", err)
		}
		n, err := writer.Write(chunk.GetChunk())
		if err != nil {
			return status.Errorf(codes.Internal, "Cannot write data to buffer %v", err)
		}
		size += n
	}
	log.Printf("Received %d bytes from the server\n", size)

	id, err := i.SaveImage(details.GetName(), data.Bytes())
	if err != nil {
		return status.Errorf(codes.Internal, "Could not save image %v", err)
	}
	res := &pb.UploadImageResponse{
		Id:   id,
		Size: uint32(size),
	}
	err = stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Unknown, "Could not send final response to the client %v", err)
	}
	return nil
}
