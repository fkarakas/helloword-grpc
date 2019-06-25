package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "grpc/helloworld/proto"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/balancer/roundrobin"
	//"google.golang.org/grpc/resolver"
)

var (
	address     = nvl(os.Getenv("SERVER"), "localhost:50051")
	defaultName = nvl(os.Getenv("POD_IP"), "<undefined>")
)

func nvl(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}

func main() {
	log.Printf("server: %s client ip: %s", address, defaultName)

	//resolver.SetDefaultScheme("dns")
	//conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			log.Printf("could not greet: %v", err)
		}
		time.Sleep(2 * time.Second)
		log.Printf("Greeting: %s", r.Message)
		cancel()
	}

}
