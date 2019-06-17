package julianclient

import (
	"context"
	"log"
	"time"

	"planetpositions/julian/grpc/v1"

	"google.golang.org/grpc"
)

type JulianClient struct {
	address string
}

func (j *JulianClient) newConnection() v1.JulianServiceClient {

	// Should I pool this?
	// Set up a connection to the server.
	conn, err := grpc.Dial(j.address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	return v1.NewJulianServiceClient(conn)
}

func (j *JulianClient) Convert(year, month, day int32, hour float64) (*v1.ConvertResponse, error) {

	c := j.newConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert
	req := v1.ConvertRequest{
		Year:  year,
		Month: month,
		Day:   day,
		Hour:  hour,
	}
	return c.Convert(ctx, &req)
}

func (j *JulianClient) TimeJulianCentury(julianDay float64) float64 {
	c := j.newConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TimeJulianCentury
	req := v1.JulianRequest{
		JulianDay: julianDay,
	}
	return c.TimeJulianCentury(ctx, &req)
}
