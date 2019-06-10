package julianclient

import (
	"context"
	"log"
	"time"

	"planetpositions/julian/grpc/v1"

	"google.golang.org/grpc"
)

type julianClient struct {
	address string
}

func (j *julianClient) Convert(year, month, day int32, hour float64) (*v1.ConvertResponse, error) {

	// Should I pool this?
	// Set up a connection to the server.
	conn, err := grpc.Dial(j.address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewJulianServiceClient(conn)

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
