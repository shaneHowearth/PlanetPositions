package main

import (
	"fmt"
	"net/http"
	"strconv"

	sun "planetpositions/sun/pkg/v1/client"

	"github.com/go-chi/chi"
)

var sc = sun.SunClient{Address: "sun.planet_positions:5055"}

func planetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/Sunrise/{long}/{lat}/{year}/{month}/{day}", GetSunrise)
	return router
}

// GetSunrise -
func GetSunrise(w http.ResponseWriter, r *http.Request) {
	long, err := strconv.ParseFloat(chi.URLParam(r, "long"), 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed longitude")
	}
	lat, err := strconv.ParseFloat(chi.URLParam(r, "lat"), 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed latitude")
	}
	year, err := strconv.Atoi(chi.URLParam(r, "year"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed year")
	}
	month, err := strconv.Atoi(chi.URLParam(r, "month"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed month")
	}
	day, err := strconv.Atoi(chi.URLParam(r, "day"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed day")
	}
	st, err := sc.GetSunrise(long, lat, int32(year), int32(month), int32(day))
	if err != nil {
		// TODO
		// log the error
		fmt.Printf("An error occurred with GetSunrise with Y: %d, M: %d, D: %d, Long: %f, Lat: %f, Error: %v", year, month, day, long, lat, err)
		respondWithError(w, http.StatusInternalServerError, "An unexpected error has occurred, the issue has been reported to our engineers and will be looked into")
	}
	respondWithJSON(w, http.StatusOK, st)

}
