package v1

import (
	"context"
	"fmt"
	"math"
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
	s.Address = "julian:5055"
	return &s
}

func (s *sunServiceServer) GetSunrise(ctx context.Context, req *v1.SunriseRequest) (*v1.SunriseTime, error) {
	// Validate input
	if ok, err := isValidInput(req.Year, req.Month, req.Day, req.Hour); !ok {
		return nil, fmt.Errorf("unusable input provided: %v", err)
	}

	// Convert the date to a Julian date
	jd, err := s.Convert(req.Year, req.Month, req.Day, req.Hour)
	if err != nil {
		return nil, err
	}
	// Calculate Sunrise/Sunset
	sunriseJD, err := s.SunriseUTC(jd.JulianDateTime, req.Latitude, req.Longitude)
	if err != nil {
		return nil, err
	}

	cal, err := s.DayFromJulianDay(sunriseJD)
	if err != nil {
		return nil, err
	}

	h := sunriseJD - math.Floor(sunriseJD)

	// Convert the julian date to y, m, d
	return &v1.SunriseTime{
		Year:  cal.Year,
		Month: cal.Month,
		Day:   cal.Day,
		Hour:  h,
	}, nil
}
