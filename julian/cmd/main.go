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
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	v1.RegisterJulianServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Convert -
func (s *server) Convert(ctx context.Context, req *v1.ConvertRequest) (*v1.JulianResponse, error) {
	jd := julian.GetJulianDay(req.GetYear(), req.GetMonth(), req.GetDay(), req.GetHour())
	return &v1.JulianResponse{JulianDateTime: jd}, nil
}

// TimeJulianCentury -
func (s *server) TimeJulianCentury(ctx context.Context, req *v1.JulianRequest) (*v1.JulianResponse, error) {
	tjc := julian.TimeJulianCentury(req.GetJulianDateTime())
	return &v1.JulianResponse{JulianDateTime: tjc}, nil
}

// JulianDayFromJulianCentury -
func (s *server) JulianDayFromJulianCentury(ctx context.Context, req *v1.JulianRequest) (*v1.JulianResponse, error) {

	tjc := julian.JulianDayFromJulianCentury(req.GetJulianDateTime())
	return &v1.JulianResponse{JulianDateTime: tjc}, nil
}

// DayFromJulianDay -
func (s *server) DayFromJulianDay(ctx context.Context, req *v1.JulianRequest) (*v1.CalendarResponse, error) {

	y, m, d := julian.DayFromJulianDay(req.GetJulianDateTime())
	return &v1.CalendarResponse{Year: y, Month: m, Day: d}, nil
}
