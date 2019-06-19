package v1

import (
	"context"
	"fmt"
	"planetpositions/sun/grpc/v1"

	jc "planetpositions/julian/pkg/v1/client"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// sunServiceServer is implementation of v1.SunServiceServer proto interface
type sunServiceServer struct {
	jc.JulianClient
}

// NewSunService creates Sun service
func NewSunService() v1.SunServiceServer {
	s := sunServiceServer{}
	s.Address = "julian"
	return &s
}

func (s *sunServiceServer) GetSunrise(ctx context.Context, req *v1.SunriseRequest) (*v1.SunriseTime, error) {
	// Validate input
	if ok, err := isValidInput(req.Year, req.Month, req.Day, req.Hour); !ok {
		return nil, fmt.Errorf("unusable input provided: %v", err)
	}
	//date := req.GetDate()
	jd, err := s.Convert(req.Year, req.Month, req.Day, req.Hour)
	if err != nil {
		return nil, err
	}
	fmt.Println(jd)
	// Calculate Sunrise/Sunset
	/*
		return &v1.SunriseTime{
			Year:  resp,
			Month: resp,
			Day:   resp,
			Hour:  resp,
		}, nil
	*/
	return nil, nil
}
