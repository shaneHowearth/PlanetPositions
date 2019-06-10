package main

import (
	"context"
	"log"
	"net"

	pb "planetpositions/julian/grpc/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterJulianServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Create(ctx context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	return &pb.ConvertResponse{}, nil
}
