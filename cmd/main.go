package main

import (
	"log"
	"net"

	pb "bloombox/api/bloom_pb"
	"bloombox/internal/server"

	"google.golang.org/grpc"
)

func main() {
	// 1. Port
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// 2. create a gRPC server
	grpcEngine := grpc.NewServer()

	// 3. init handler
	bloomLogic := &server.BloomServer{}

	// 4. register handler with the enginer
	pb.RegisterBloomServiceServer(grpcEngine, bloomLogic)

	// 5. Start serving
	log.Printf("BloomBox gRPC server listening at port %v", lis.Addr())
	if err := grpcEngine.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
