package v1

import (
	"context"
	"database/sql"
	"fmt"
	"planetpositions/sun/grpc/v1"

	jc "planetpositions/julian/pkg/v1/client"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// sunServiceServer is implementation of v1.SunServiceServer proto interface
type sunServiceServer struct {
	jc.JulianClient
	db *sql.DB
}

// NewSunService creates Sun service
func NewSunService(db *sql.DB) v1.SunServiceServer {
	return &sunServiceServer{db: db}
}

func isValidInput(year, month, day int32, hour float64) (bool, error) {
	if day <= 0 {
		return false, fmt.Errorf("invalid day supplied")
	}
	if month <= 0 || month > 12 {
		return false, fmt.Errorf("invalid month supplied")
	}
	thirtyOnes := map[int32]string{
		1:  "January",
		3:  "March",
		5:  "May",
		7:  "July",
		8:  "August",
		10: "October",
		12: "December",
	}
	if _, ok := thirtyOnes[month]; ok {
		if day > 31 {
			return false, fmt.Errorf("there are only 31 days in %s", thirtyOnes[month])
		}
	}
	thirtys := map[int32]string{
		4:  "April",
		6:  "June",
		9:  "September",
		11: "November",
	}
	if _, ok := thirtys[month]; ok {
		if day > 30 {
			return false, fmt.Errorf("there are only 30 days in %s", thirtys[month])
		}
	}
	if month == 2 {
		// Leap Year calculation
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
			if day > 29 {
				return false, fmt.Errorf("there are only 29 days in February during leap years")
			}
		} else {
			if day > 28 {
				return false, fmt.Errorf("there are only 28 days in February during non-leap years")

			}
		}
	}
	if year < -1000 || year > 3000 {
		return false, fmt.Errorf("the algorithm used is not valid for years outside of the range -1000 to 3000")
	}
	return true, nil
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
	// Calculate SUnrise/Sunset
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

// checkAPI checks if the API version requested by client is supported by server
func (s *sunServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}
