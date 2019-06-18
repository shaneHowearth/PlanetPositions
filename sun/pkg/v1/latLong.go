package v1

import (
	"fmt"
	"math"
)

func calcDayOfYear(month, day int32, leapYear bool) int32 {
	k := int32(1)
	if !leapYear {
		k = 2
	}
	dayOfYear := int32(math.Floor(float64(275*month))/float64(9)) - k*int32(math.Floor(float64(month+9)/float64(12))) + day - int32(30)
	return dayOfYear
}

func (s *sunServiceServer) geometricMeanLongitudeSun(t float64) float64 {
	l0 := 280.46646 + t*(36000.76983+0.0003032*t)
	l0 = math.Mod(l0, 360)
	for {
		if l0 < 360 {
			l0 += 360
		} else {
			break
		}
	}
	return l0 // In Degrees
}

func (s *sunServiceServer) geometricMeanAnamolySun(t float64) float64 {
	// t is the number of Julian centuries since J2000.0
	return 357.52911 + t*(35999.05029-0.0001537*t)
}

func (s *sunServiceServer) eccentricityEarthOrbit(t float64) float64 {
	return 0.016708634 - t*(0.000042037+0.0000001267*t)
}

func (s *sunServiceServer) equationCentreSun(t float64) float64 {
	m := s.geometricMeanAnamolySun(t)
	mrad := degreesToRadians(m)
	sinm := math.Sin(mrad)
	sin2m := math.Sin(mrad + mrad)
	sin3m := math.Sin(mrad + mrad + mrad)

	return sinm*(1.914602-t*(0.004817+0.000014*t)) + sin2m*(0.019993-0.000101*t) + sin3m*0.000289 // In Degrees
}

func (s *sunServiceServer) sunTrueLongitude(t float64) float64 {
	l0 := s.geometricMeanLongitudeSun(t)
	c := s.equationCentreSun(t)

	return l0 + c // In Degrees
}

func (s *sunServiceServer) sunTrueAnamoly(t float64) float64 {
	m := s.geometricMeanAnamolySun(t)
	c := s.equationCentreSun(t)

	return m + c // In Degrees

}

func (s *sunServiceServer) sunRadiusVector(t float64) float64 {
	v := s.sunTrueAnamoly(t)
	e := s.eccentricityEarthOrbit(t)

	return (1.000001018 * (1 - e*e)) / (1 + e*math.Cos(degreesToRadians(v))) // In Astronomical Units
}

func (s *sunServiceServer) sunApparentLongitude(t float64) float64 {
	o := s.sunTrueLongitude(t)

	omega := 125.04 - 1934.136*t
	return o - 0.00569 - 0.00478*math.Sin(degreesToRadians(omega)) // In Degrees
}

func (s *sunServiceServer) meanObliquityOfEcliptic(t float64) float64 {

	seconds := 21.448 - t*(46.8150+t*(0.00059-t*(0.001813)))
	return 23.0 + (26.0+(seconds/60.0))/60.0 // In Degrees
}

func (s *sunServiceServer) obliquityCorrection(t float64) float64 {

	e0 := s.meanObliquityOfEcliptic(t)

	omega := 125.04 - 1934.136*t
	return e0 + 0.00256*math.Cos(degreesToRadians(omega))
}

// SunRightAscension -
func (s *sunServiceServer) SunRightAscension(t float64) float64 {
	e := s.obliquityCorrection(t)
	lambda := s.sunApparentLongitude(t)

	tananum := (math.Cos(degreesToRadians(e)) * math.Sin(degreesToRadians(lambda)))
	tanadenom := (math.Cos(degreesToRadians(lambda)))
	//TODO TEST THIS BECAUSE Atan2 might not be the same
	return radiansToDegrees(math.Atan2(tananum, tanadenom)) // In Degrees
}
func (s *sunServiceServer) sunDeclination(t float64) float64 {
	e := s.obliquityCorrection(t)
	lambda := s.sunApparentLongitude(t)

	sint := math.Sin(degreesToRadians(e)) * math.Sin(degreesToRadians(lambda))
	return radiansToDegrees(math.Asin(sint)) // In Degrees
}
func (s *sunServiceServer) equationOfTime(t float64) float64 {
	epsilon := s.obliquityCorrection(t)
	l0 := s.geometricMeanLongitudeSun(t)
	e := s.eccentricityEarthOrbit(t)
	m := s.geometricMeanAnamolySun(t)

	y := math.Tan(degreesToRadians(epsilon) / 2.0)
	y *= y

	sin2l0 := math.Sin(2.0 * degreesToRadians(l0))
	sinm := math.Sin(degreesToRadians(m))
	cos2l0 := math.Cos(2.0 * degreesToRadians(l0))
	sin4l0 := math.Sin(4.0 * degreesToRadians(l0))
	sin2m := math.Sin(2.0 * degreesToRadians(m))

	Etime := y*sin2l0 - 2.0*e*sinm + 4.0*e*y*sinm*cos2l0 - 0.5*y*y*sin4l0 - 1.25*e*e*sin2m

	return radiansToDegrees(Etime) * 4.0 // In minutes of time
}

func hourAngle(lat, solarDeclination float64) float64 {
	latRad := degreesToRadians(lat)
	sdRad := degreesToRadians(solarDeclination)

	// HAarg := (math.Cos(degreesToRadians(90.833))/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad))

	HA := (math.Acos(math.Cos(degreesToRadians(90.833))/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad)))

	return HA // In Radians
}

func (s *sunServiceServer) HourAngleSunrise(lat, solarDec float64) float64 {
	return hourAngle(lat, solarDec) // in radians
}

func (s *sunServiceServer) HourAngleSunset(lat, solarDec float64) float64 {
	// Negate the hour angle for sunset
	return -hourAngle(lat, solarDec) // in radians
}

func (s *sunServiceServer) solNoonUTC(t, longitude float64) (float64, error) {
	// First pass uses approximate solar noon to calculate eqtime
	jd, err := s.JulianDayFromJulianCentury(t)
	if err != nil {
		return 0, fmt.Errorf("Solnoon encountered the following error when executing first JulianDayFromJulianCentury: %v", err)
	}
	tnoon, err := s.TimeJulianCentury(jd.JulianDateTime + longitude/360.0)
	if err != nil {
		return 0, fmt.Errorf("Solnoon encountered the following error when executing TimeJulianCentury for tnoon: %v", err)
	}
	eqTime := s.equationOfTime(tnoon.JulianDateTime)
	solNoonUTC := 720 + (longitude * 4) - eqTime // min

	jd, err = s.JulianDayFromJulianCentury(t)
	if err != nil {
		return 0, fmt.Errorf("Solnoon encountered the following error when executing second JulianDayFromJulianCentury: %v", err)
	}
	newt, err := s.TimeJulianCentury(jd.JulianDateTime - 0.5 + solNoonUTC/1440.0)
	if err != nil {
		return 0, fmt.Errorf("Solnoon encountered the following error when executing TimeJulianCentury for newt: %v", err)
	}

	eqTime = s.equationOfTime(newt.JulianDateTime)
	// var solarNoonDec = calcSunDeclination(newt)
	solNoonUTC = 720 + (longitude * 4) - eqTime // min

	return solNoonUTC, nil
}

func (s *sunServiceServer) sunriseUTC(JD, latitude, longitude float64) (float64, error) {
	t, err := s.TimeJulianCentury(JD)
	if err != nil {
		return 0, fmt.Errorf("sunriseUTC encountered the following error when executing TimeJulianCentury for t: %v", err)
	}

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin, err := s.solNoonUTC(t.JulianDateTime, longitude)
	if err != nil {
		return 0, fmt.Errorf("sunriseUTC encountered the following error when executing solNoonUTC for noonmin: %v", err)
	}
	tnoon, err := s.TimeJulianCentury(JD + noonmin/1440.0)
	if err != nil {
		return 0, fmt.Errorf("sunriseUTC encountered the following error when executing TimeJulianCentury for tnoon: %v", err)
	}

	// *** First pass to approximate sunrise (using solar noon)

	eqTime := s.equationOfTime(tnoon.JulianDateTime)
	solarDec := s.sunDeclination(tnoon.JulianDateTime)
	hourAngle := s.HourAngleSunrise(latitude, solarDec)

	delta := longitude - radiansToDegrees(hourAngle)
	timeDiff := 4 * delta              // in minutes of time
	timeUTC := 720 + timeDiff - eqTime // in minutes

	// *** Second pass includes fractional jday in gamma calc
	jd, err := s.JulianDayFromJulianCentury(t.JulianDateTime)
	if err != nil {
		return 0, fmt.Errorf("sunriseUTC encountered the following error when executing JulianDayFromJulianCentury for jd: %v", err)
	}
	newt, err := s.TimeJulianCentury(jd.JulianDateTime + timeUTC/1440.0)
	if err != nil {
		return 0, fmt.Errorf("sunriseUTC encountered the following error when executing TimeJulianCentury for newt: %v", err)
	}
	eqTime = s.equationOfTime(newt.JulianDateTime)
	solarDec = s.sunDeclination(newt.JulianDateTime)
	hourAngle = s.HourAngleSunrise(latitude, solarDec)
	delta = longitude - radiansToDegrees(hourAngle)
	timeDiff = 4 * delta
	timeUTC = 720 + timeDiff - eqTime // in minutes

	return timeUTC, nil
}

func (s *sunServiceServer) sunsetUTC(JD, latitude, longitude float64) (float64, error) {
	t, err := s.TimeJulianCentury(JD)
	if err != nil {
		return 0, fmt.Errorf("sunsetUTC encountered the following error when executing TimeJulianCentury for t: %v", err)
	}

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin, err := s.solNoonUTC(t.JulianDateTime, longitude)
	if err != nil {
		return 0, fmt.Errorf("sunsetUTC encountered the following error when executing solNoonUTC for noonmin: %v", err)
	}
	tnoon, err := s.TimeJulianCentury(JD + noonmin/1440.0)
	if err != nil {
		return 0, fmt.Errorf("sunsetUTC encountered the following error when executing TimeJulianCentury for tnoon: %v", err)
	}

	// First calculates sunrise and approx length of day

	eqTime := s.equationOfTime(tnoon.JulianDateTime)
	solarDec := s.sunDeclination(tnoon.JulianDateTime)
	hourAngle := s.HourAngleSunset(latitude, solarDec)

	delta := longitude - radiansToDegrees(hourAngle)
	timeDiff := 4 * delta
	timeUTC := 720 + timeDiff - eqTime

	// first pass used to include fractional day in gamma calc
	jd, err := s.JulianDayFromJulianCentury(t.JulianDateTime)
	if err != nil {
		return 0, fmt.Errorf("sunsetUTC encountered the following error when executing JulianDayFromJulianCentury for jd: %v", err)
	}
	newt, err := s.TimeJulianCentury(jd.JulianDateTime + timeUTC/1440.0)
	if err != nil {
		return 0, fmt.Errorf("sunsetUTC encountered the following error when executing TimeJulianCentury for newt: %v", err)
	}
	eqTime = s.equationOfTime(newt.JulianDateTime)
	solarDec = s.sunDeclination(newt.JulianDateTime)
	hourAngle = s.HourAngleSunset(latitude, solarDec)

	delta = longitude - radiansToDegrees(hourAngle)
	timeDiff = 4 * delta
	timeUTC = 720 + timeDiff - eqTime // in minutes

	return timeUTC, nil
}
