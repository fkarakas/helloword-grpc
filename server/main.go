package main

import (
	"context"
	"log"
	"net"
	"os"
	//"time"

	pb "grpc/helloworld/proto"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

var (
	port        = ":50051"
	defaultName = nvl(os.Getenv("POD_IP"), "<undefined>")
)

func nvl(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received message from: %v", in.Name)
	return &pb.HelloReply{Message: "Hello from: " + defaultName}, nil
}

func main() {
	log.Printf("Listening on %s", port)

	//keepAliveMaxConnectionAge := time.Second * 5

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer( /*
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: keepAliveMaxConnectionAge,
		}),*/
	)

	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
