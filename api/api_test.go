package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{}

// Email service does nothing at the moment.
func (s *server) Email(ctx context.Context, in *EmailRequest) (*EmailResponse, error) {
	return nil, nil
}

// InferImage receives streaming input and build output sending it back on the stream.
func (s *server) InferImage(stream Api_InferImageServer) error {
	request := new(InferImageRequest)
	request.Images = make([]*Image, 0, 0)

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		request.LabelPath = in.LabelPath
		request.ModelPath = in.ModelPath
		request.Images = append(request.Images, in.Images...)
	}

	response := new(InferImageResponse)
	response.Outputs = make([]*InferOutput, len(request.Images))
	for i := range response.Outputs {
		response.Outputs[i] = new(InferOutput)
		response.Outputs[i].Name = request.Images[i].Name
		response.Outputs[i].Label = "all ok"
	}

	return stream.SendAndClose(response)
}

// newServer starts a new gRPC servers and listens on the port
func newServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterApiServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func TestApiClient_InferImage(t *testing.T) {
	go newServer()

	// Set up a connection to the newServer.
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := NewApiClient(conn)

	stream, err := client.InferImage(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// send one image at a time via stream
	for i := 0; i < 3; i++ {
		in := new(InferImageRequest)
		in.Images = make([]*Image, 1)
		in.ModelPath = "modelPath"
		in.LabelPath = "labelPath"

		in.Images[0] = new(Image)
		in.Images[0].Name = fmt.Sprintf("image-%d", i)
		in.Images[0].Data = []byte{byte(i), 0, 1, 2}

		if err := stream.Send(in); err != nil {
			t.Fatal(err)
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
	}

	if len(reply.Outputs) != 3 {
		t.Fatalf("expected output len to be %d, got %d", 3, len(reply.Outputs))
	}
}
