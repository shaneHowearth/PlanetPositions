package v1

import (
	"context"
	"database/sql"
	"planetpositions/sun/grpc/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// sunServiceServer is implementation of v1.SunServiceServer proto interface
type sunServiceServer struct {
	db *sql.DB
}

// NewSunServiceServer creates Sun service
func NewSunService(db *sql.DB) v1.SunServiceServer {
	return &sunServiceServer{db: db}
}

func (s *sunServiceServer) GetSunrise(ctx context.Context, req *v1.SunriseRequest) (*v1.SunriseTime, error) {
	return &v1.SunriseTime{}, nil
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
