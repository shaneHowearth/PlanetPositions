package main

import (
	"context"
	"log"
	"net"

	v1 "planetpositions/sun/grpc/v1"
	sun "planetpositions/sun/pkg/v1/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

const (
	port = ":50051"
)

var ss = sun.NewSunService()

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	v1.RegisterSunServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// GetSunrise -
func (s *server) GetSunrise(ctx context.Context, req *v1.SunriseRequest) (*v1.SunriseTime, error) {
	st, err := ss.GetSunrise(ctx, req)
	if err != nil {
		return nil, err
	}
	return st, nil
}
