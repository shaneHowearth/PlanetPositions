package sunclient

import (
	"context"
	"log"
	"time"

	"planetpositions/sun/grpc/v1"

	"google.golang.org/grpc"
)

// SunClient -
type SunClient struct {
	Address string
}

func (s *SunClient) newConnection() v1.SunServiceClient {

	// Set up a connection to the server.
	conn, err := grpc.Dial(s.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	return v1.NewSunServiceClient(conn)
}

// GetSunrise -
func (s *SunClient) GetSunrise(long, lat float64, year, month, day int32) (*v1.SunriseTime, error) {

	c := s.newConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := v1.SunriseRequest{
		Api:       "v1",
		Longitude: long,
		Latitude:  lat,
		Year:      year,
		Month:     month,
		Day:       day,
	}
	return c.GetSunrise(ctx, &req)
}
