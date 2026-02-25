package server

import (
	"context"
	"fmt"

	"sync"

	pb "bloombox/api/bloom_pb"
	"bloombox/internal/bloom"
)

type BloomServer struct {
	pb.UnimplementedBloomServiceServer

	// actual backend logic
	activeFilter bloom.Filter
	mu           sync.RWMutex // protect racing condition e.g.: when client A trying to check but client B trying to replace active Filter
}

// create filter will handle incoming gRPC requests --> init a new filter
func (s *BloomServer) CreateFilter(ctx context.Context, req *pb.CreateFilterRequest) (*pb.CreateFilterResponse, error) {
	cfg := bloom.Config{
		Type:          req.FilterType,
		ExpectedItems: int(req.ExpectedItems),
		FalsePosRate:  req.FalsePosRate,
		GrowthFactor:  req.GrowthFactor,
	}

	filter, err := bloom.NewFilter(cfg)
	if err != nil {
		return &pb.CreateFilterResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	s.mu.Lock()
	s.activeFilter = filter
	s.mu.RUnlock()
	return &pb.CreateFilterResponse{
		Success: true,
		Message: fmt.Sprintf("%s filter created successfully", req.FilterType),
	}, nil
}

// Add handler
func (s *BloomServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	s.mu.RLock()
	filter := s.activeFilter // we grab the pointer to that filter, even activeFilter got replaced, does not crash
	s.mu.RUnlock()
	if filter == nil {
		return &pb.AddResponse{Success: false}, nil
	}

	filter.Add(req.Data)
	return &pb.AddResponse{Success: true}, nil
}

// Check handler
func (s *BloomServer) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	s.mu.RLock()
	filter := s.activeFilter // we grab the pointer to that filter, even activeFilter got replaced, does not crash
	s.mu.RUnlock()

	if filter == nil {
		return &pb.CheckResponse{Found: false}, nil
	}
	found := filter.Check(req.Data)
	return &pb.CheckResponse{Found: found}, nil

}
