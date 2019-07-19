package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "grpc/helloworld/proto"

	"github.com/sercand/kuberesolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
)

var (
	address     = nvl(os.Getenv("SERVER"), "localhost:50051")
	lb          = nvl(os.Getenv("LB"), "")
	defaultName = nvl(os.Getenv("POD_IP"), "<undefined>")
)

func nvl(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func DialWithLoadBalancer(lb string, addr string) (*grpc.ClientConn, error) {
	switch lb {
	case "dns":
		resolver.SetDefaultScheme("dns")
		return grpc.Dial(address, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	case "kube":
		kuberesolver.RegisterInCluster()
		//return grpc.Dial("kubernetes:///service-name.namespace:portname", grpc.WithInsecure())
		return grpc.Dial("kubernetes:///"+address, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	default:
		return grpc.Dial(address, grpc.WithInsecure())
	}
}

func main() {
	log.Printf("server: %s client ip: %s", address, defaultName)

	conn, err := DialWithLoadBalancer(lb, address)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
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
