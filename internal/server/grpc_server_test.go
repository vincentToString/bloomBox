package server

import (
	"context"
	"net"
	"testing"

	pb "bloombox/api/bloom_pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

// in memory gRPC server
func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterBloomServiceServer(s, &BloomServer{})

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

// bufDialer connects to in memory server
func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestBloomServer(t *testing.T) {
	ctx := context.Background()

	// Set up the client
	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewBloomServiceClient(conn)

	// 1. a standard filter
	createRes, err := client.CreateFilter(ctx, &pb.CreateFilterRequest{
		FilterType:    "standard",
		ExpectedItems: 100,
		FalsePosRate:  0.01,
	})
	if err != nil || !createRes.Success {
		t.Fatalf("Failed to create filter: %v", err)
	}

	// 2. Check add function
	addRes, err := client.Add(ctx, &pb.AddRequest{Data: []byte("grpc_test_item")})
	if err != nil || !addRes.Success {
		t.Fatalf("Failed to add item: %v", err)
	}

	// 3. Check function
	checkRes, err := client.Check(ctx, &pb.CheckRequest{Data: []byte("grpc_test_item")})
	if err != nil || !checkRes.Found {
		t.Errorf("Expected to find 'grpc_test_item' but it was not found with err: %v", err)
	}

	// 4. Check function (non-existed item)
	checkResMissing, err := client.Check(ctx, &pb.CheckRequest{Data: []byte("does_not_exist")})
	if err != nil || checkResMissing.Found {
		t.Errorf("Did not expect to find un-added item; error: %v", err)
	}
}
