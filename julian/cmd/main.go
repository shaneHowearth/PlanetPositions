package main

import (
	"context"
	"log"
	"net"

	v1 "planetpositions/julian/grpc/v1"
	julian "planetpositions/julian/pkg/v1/service"

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
	v1.RegisterJulianServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Create -
func (s *server) Create(ctx context.Context, req *v1.ConvertRequest) (*v1.ConvertResponse, error) {
	jd := julian.GetJulianDay(req.GetYear(), req.GetMonth(), req.GetDay(), req.GetHour())
	return &v1.ConvertResponse{JulianDateTime: jd}, nil
}
