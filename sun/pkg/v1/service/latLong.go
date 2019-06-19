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

// GeometricMeanLongitudeSun -
func (s *sunServiceServer) GeometricMeanLongitudeSun(t float64) float64 {
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

// GeometricMeanAnamolySun -
func (s *sunServiceServer) GeometricMeanAnamolySun(t float64) float64 {
	// t is the number of Julian centuries since J2000.0
	return 357.52911 + t*(35999.05029-0.0001537*t)
}

// EccentricityEarthOrbit -
func (s *sunServiceServer) EccentricityEarthOrbit(t float64) float64 {
	return 0.016708634 - t*(0.000042037+0.0000001267*t)
}

// EquationCentreSun -
func (s *sunServiceServer) EquationCentreSun(t float64) float64 {
	m := s.GeometricMeanAnamolySun(t)
	mrad := degreesToRadians(m)
	sinm := math.Sin(mrad)
	sin2m := math.Sin(mrad + mrad)
	sin3m := math.Sin(mrad + mrad + mrad)

	return sinm*(1.914602-t*(0.004817+0.000014*t)) + sin2m*(0.019993-0.000101*t) + sin3m*0.000289 // In Degrees
}

// SunTrueLongitude -
func (s *sunServiceServer) SunTrueLongitude(t float64) float64 {
	l0 := s.GeometricMeanLongitudeSun(t)
	c := s.EquationCentreSun(t)

	return l0 + c // In Degrees
}

// SunTrueAnamoly -
func (s *sunServiceServer) SunTrueAnamoly(t float64) float64 {
	m := s.GeometricMeanAnamolySun(t)
	c := s.EquationCentreSun(t)

	return m + c // In Degrees

}

// SunRadiusVector -
func (s *sunServiceServer) SunRadiusVector(t float64) float64 {
	v := s.SunTrueAnamoly(t)
	e := s.EccentricityEarthOrbit(t)

	return (1.000001018 * (1 - e*e)) / (1 + e*math.Cos(degreesToRadians(v))) // In Astronomical Units
}

// SunApparentLongitude -
func (s *sunServiceServer) SunApparentLongitude(t float64) float64 {
	o := s.SunTrueLongitude(t)

	omega := 125.04 - 1934.136*t
	return o - 0.00569 - 0.00478*math.Sin(degreesToRadians(omega)) // In Degrees
}

// MeanObliquityOfEcliptic -
func (s *sunServiceServer) MeanObliquityOfEcliptic(t float64) float64 {

	seconds := 21.448 - t*(46.8150+t*(0.00059-t*(0.001813)))
	return 23.0 + (26.0+(seconds/60.0))/60.0 // In Degrees
}

// ObliquityCorrection -
func (s *sunServiceServer) ObliquityCorrection(t float64) float64 {

	e0 := s.MeanObliquityOfEcliptic(t)

	omega := 125.04 - 1934.136*t
	return e0 + 0.00256*math.Cos(degreesToRadians(omega))
}

// SunRightAscension -
func (s *sunServiceServer) SunRightAscension(t float64) float64 {
	e := s.ObliquityCorrection(t)
	lambda := s.SunApparentLongitude(t)

	tananum := (math.Cos(degreesToRadians(e)) * math.Sin(degreesToRadians(lambda)))
	tanadenom := (math.Cos(degreesToRadians(lambda)))
	//TODO TEST THIS BECAUSE Atan2 might not be the same
	return radiansToDegrees(math.Atan2(tananum, tanadenom)) // In Degrees
}

// SunDeclination -
func (s *sunServiceServer) SunDeclination(t float64) float64 {
	e := s.ObliquityCorrection(t)
	lambda := s.SunApparentLongitude(t)

	sint := math.Sin(degreesToRadians(e)) * math.Sin(degreesToRadians(lambda))
	return radiansToDegrees(math.Asin(sint)) // In Degrees
}

// EquationOfTime -
func (s *sunServiceServer) EquationOfTime(t float64) float64 {
	epsilon := s.ObliquityCorrection(t)
	l0 := s.GeometricMeanLongitudeSun(t)
	e := s.EccentricityEarthOrbit(t)
	m := s.GeometricMeanAnamolySun(t)

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

// HourAngleSunrise -
func (s *sunServiceServer) HourAngleSunrise(lat, solarDec float64) float64 {
	return hourAngle(lat, solarDec) // in radians
}

// HourAngleSunset -
func (s *sunServiceServer) HourAngleSunset(lat, solarDec float64) float64 {
	// Negate the hour angle for sunset
	return -hourAngle(lat, solarDec) // in radians
}

// SolNoonUTC -
func (s *sunServiceServer) SolNoonUTC(t, longitude float64) (float64, error) {
	// First pass uses approximate solar noon to calculate eqtime
	jd, err := s.JulianDayFromJulianCentury(t)
	if err != nil {
		return 0, fmt.Errorf("solnoon encountered the following error when executing first JulianDayFromJulianCentury: %v", err)
	}
	tnoon, err := s.TimeJulianCentury(jd.JulianDateTime + longitude/360.0)
	if err != nil {
		return 0, fmt.Errorf("solnoon encountered the following error when executing TimeJulianCentury for tnoon: %v", err)
	}
	eqTime := s.EquationOfTime(tnoon.JulianDateTime)
	solNoonUTC := 720 + (longitude * 4) - eqTime // min

	jd, err = s.JulianDayFromJulianCentury(t)
	if err != nil {
		return 0, fmt.Errorf("solnoon encountered the following error when executing second JulianDayFromJulianCentury: %v", err)
	}
	newt, err := s.TimeJulianCentury(jd.JulianDateTime - 0.5 + solNoonUTC/1440.0)
	if err != nil {
		return 0, fmt.Errorf("solnoon encountered the following error when executing TimeJulianCentury for newt: %v", err)
	}

	eqTime = s.EquationOfTime(newt.JulianDateTime)
	// var solarNoonDec = calcSunDeclination(newt)
	solNoonUTC = 720 + (longitude * 4) - eqTime // min

	return solNoonUTC, nil
}

// SunriseUTC -
func (s *sunServiceServer) SunriseUTC(JD, latitude, longitude float64) (float64, error) {
	t, err := s.TimeJulianCentury(JD)
	if err != nil {
		return 0, fmt.Errorf("sunriseUTC encountered the following error when executing TimeJulianCentury for t: %v", err)
	}

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin, err := s.SolNoonUTC(t.JulianDateTime, longitude)
	if err != nil {
		return 0, fmt.Errorf("sunriseUTC encountered the following error when executing solNoonUTC for noonmin: %v", err)
	}
	tnoon, err := s.TimeJulianCentury(JD + noonmin/1440.0)
	if err != nil {
		return 0, fmt.Errorf("sunriseUTC encountered the following error when executing TimeJulianCentury for tnoon: %v", err)
	}

	// *** First pass to approximate sunrise (using solar noon)

	eqTime := s.EquationOfTime(tnoon.JulianDateTime)
	solarDec := s.SunDeclination(tnoon.JulianDateTime)
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
	eqTime = s.EquationOfTime(newt.JulianDateTime)
	solarDec = s.SunDeclination(newt.JulianDateTime)
	hourAngle = s.HourAngleSunrise(latitude, solarDec)
	delta = longitude - radiansToDegrees(hourAngle)
	timeDiff = 4 * delta
	timeUTC = 720 + timeDiff - eqTime // in minutes

	return timeUTC, nil
}

// SunsetUTC -
func (s *sunServiceServer) SunsetUTC(JD, latitude, longitude float64) (float64, error) {
	t, err := s.TimeJulianCentury(JD)
	if err != nil {
		return 0, fmt.Errorf("sunsetUTC encountered the following error when executing TimeJulianCentury for t: %v", err)
	}

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin, err := s.SolNoonUTC(t.JulianDateTime, longitude)
	if err != nil {
		return 0, fmt.Errorf("sunsetUTC encountered the following error when executing solNoonUTC for noonmin: %v", err)
	}
	tnoon, err := s.TimeJulianCentury(JD + noonmin/1440.0)
	if err != nil {
		return 0, fmt.Errorf("sunsetUTC encountered the following error when executing TimeJulianCentury for tnoon: %v", err)
	}

	// First calculates sunrise and approx length of day

	eqTime := s.EquationOfTime(tnoon.JulianDateTime)
	solarDec := s.SunDeclination(tnoon.JulianDateTime)
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
	eqTime = s.EquationOfTime(newt.JulianDateTime)
	solarDec = s.SunDeclination(newt.JulianDateTime)
	hourAngle = s.HourAngleSunset(latitude, solarDec)

	delta = longitude - radiansToDegrees(hourAngle)
	timeDiff = 4 * delta
	timeUTC = 720 + timeDiff - eqTime // in minutes

	return timeUTC, nil
}

/*
func findRecentSunrise(julianDay, latitude, longitude float64) float64 {

	time := sunriseUTC(julianDay, latitude, longitude)
	// ??
	/*
		while(!isNumber(time)){
			julianDay -= 1.0
			time = sunriseUTC(julianDay, latitude, longitude)
		}
*/
/*
	return julianDay
}
*/

/*
func findRecentSunset(julianDay, latitude, longitude float64) float64 {
	time := sunsetUTC(julianDay, latitude, longitude)
	// TODO
	/*
		while(!isNumber(time)){
			julianDay -= 1.0
			time = sunsetUTC(julianDay, latitude, longitude)
		}
*/

/*
	return julianDay
}
*/

/*
func findNextSunrise(julianDay, latitude, longitude float64) float64 {
	time := sunriseUTC(julianDay, latitude, longitude)
	// TODO
	/*
		while(!isNumber(time)){
			julianDay += 1.0
			time = sunriseUTC(julianDay, latitude, longitude)
		}
*/

/*
	return julianDay
}
*/
/*
func findNextSunset(julianDay, latitude, longitude float64) float64 {

	time := sunsetUTC(julianDay, latitude, longitude)
	//TODO
	/*
		while(!isNumber(time)){
			julianDay += 1.0
			time = sunsetUTC(julianDay, latitude, longitude)
		}
*/

/*
	return julianDay
}
*/

/*
func sun(year, month, day int32, hour, latitude, longitude, index float64) {
	if (latitude >= -90) && (latitude < -89) {
		//alert("All latitudes between 89 and 90 S\n will be set to -89")
		latitude = -89
	}
	if (latitude <= 90) && (latitude > 89) {
		//alert("All latitudes between 89 and 90 N\n will be set to 89")
		latitude = 89
	}

	//*****	Calculate the time of sunrise

	//*********************************************************************/
//****************   NEW STUFF   ******   January, 2001   ****************
//*********************************************************************/
/*
	// get julian day from gRPC service
	JD := julianDay(parseFloat(riseSetForm["year"].value), indexRS+1, parseFloat(riseSetForm["day"].value))
	dow := calcDayOfWeek(JD)
	doy := calcDayOfYear(indexRS+1, parseFloat(riseSetForm["day"].value), isLeapYear(riseSetForm["year"].value))
	T := calcTimeJulianCent(JD)

	alpha := calcSunRtAscension(T)
	theta := calcSunDeclination(T)
	Etime := calcEquationOfTime(T)

	//riseSetForm["dbug"].value = doy

	//*********************************************************************/
/*

	eqTime := Etime
	solarDec := theta

	// Calculate sunrise for this date
	// if no sunrise is found, set flag nosunrise

	nosunrise := false

	riseTimeGMT := calcSunriseUTC(JD, latitude, longitude)
	if !isNumber(riseTimeGMT) {
		nosunrise = true
	}

	// Calculate sunset for this date
	// if no sunset is found, set flag nosunset

	nosunset = false
	setTimeGMT = calcSunsetUTC(JD, latitude, longitude)
	if !isNumber(setTimeGMT) {
		nosunset = true
	}

	daySavings = YesNo[index].value // = 0 (no) or 60 (yes)
	zone = latLongForm["hrsToGMT"].value
	if zone > 12 || zone < -12.5 {
		alert("The offset must be between -12.5 and 12.  \n Setting \"Off-Set\"=0")
		zone = "0"
		latLongForm["hrsToGMT"].value = zone
	}

	if !nosunrise {
		// Sunrise was found
		riseTimeLST = riseTimeGMT - (60 * zone) + daySavings
		//	in minutes
		riseStr = timeStringShortAMPM(riseTimeLST, JD)
		utcRiseStr = timeStringDate(riseTimeGMT, JD)

		riseSetForm["sunrise"].value = riseStr
		riseSetForm["utcsunrise"].value = utcRiseStr
	}

	if !nosunset {
		// Sunset was found
		setTimeLST = setTimeGMT - (60 * zone) + daySavings
		setStr = timeStringShortAMPM(setTimeLST, JD)
		utcSetStr = timeStringDate(setTimeGMT, JD)

		riseSetForm["sunset"].value = setStr
		riseSetForm["utcsunset"].value = utcSetStr
	}

	// Calculate solar noon for this date

	solNoonGMT = calcSolNoonUTC(T, longitude)
	solNoonLST = solNoonGMT - (60 * zone) + daySavings

	solnStr = timeString(solNoonLST)
	utcSolnStr = timeString(solNoonGMT)

	riseSetForm["solnoon"].value = solnStr
	riseSetForm["utcsolnoon"].value = utcSolnStr

	tsnoon = calcTimeJulianCent(calcJDFromJulianCent(T) - 0.5 + solNoonGMT/1440.0)

	eqTime = calcEquationOfTime(tsnoon)
	solarDec = calcSunDeclination(tsnoon)

	riseSetForm["eqTime"].value = (Math.floor(100 * eqTime)) / 100
	riseSetForm["solarDec"].value = (Math.floor(100 * (solarDec))) / 100

	//***********Convert lat and long to standard format
	convLatLong(latLongForm)

	// report special cases of no sunrise

	if nosunrise {
		riseSetForm["utcsunrise"].value = ""
		// if Northern hemisphere and spring or summer, OR
		// if Southern hemisphere and fall or winter, use
		// previous sunrise and next sunset

		if ((latitude > 66.4) && (doy > 79) && (doy < 267)) || ((latitude < -66.4) && ((doy < 83) || (doy > 263))) {
			newjd = findRecentSunrise(JD, latitude, longitude)
			newtime = calcSunriseUTC(newjd, latitude, longitude)
			-(60 * zone) + daySavings
			if newtime > 1440 {
				newtime -= 1440
				newjd += 1.0
			}
			if newtime < 0 {
				newtime += 1440
				newjd -= 1.0
			}
			riseSetForm["sunrise"].value =
				timeStringAMPMDate(newtime, newjd)
			riseSetForm["utcsunrise"].value = "prior sunrise"

			// if Northern hemisphere and fall or winter, OR
			// if Southern hemisphere and spring or summer, use
			// next sunrise and previous sunset

		} else if ((latitude > 66.4) && ((doy < 83) || (doy > 263))) || ((latitude < -66.4) && (doy > 79) && (doy < 267)) {
			newjd = findNextSunrise(JD, latitude, longitude)
			newtime = calcSunriseUTC(newjd, latitude, longitude)
			-(60 * zone) + daySavings
			if newtime > 1440 {
				newtime -= 1440
				newjd += 1.0
			}
			if newtime < 0 {
				newtime += 1440
				newjd -= 1.0
			}
			riseSetForm["sunrise"].value =
				timeStringAMPMDate(newtime, newjd)
				//					riseSetForm["sunrise"].value = calcDayFromJD(newjd)
				//						+ " " + timeStringDate(newtime, newjd)
			riseSetForm["utcsunrise"].value = "next sunrise"
		} else {
			alert("Cannot Find Sunrise!")
		}

		// alert("Last Sunrise was on day " + findRecentSunrise(JD, latitude, longitude))
		// alert("Next Sunrise will be on day " + findNextSunrise(JD, latitude, longitude))

	}

	if nosunset {
		riseSetForm["utcsunset"].value = ""
		// if Northern hemisphere and spring or summer, OR
		// if Southern hemisphere and fall or winter, use
		// previous sunrise and next sunset

		if ((latitude > 66.4) && (doy > 79) && (doy < 267)) || ((latitude < -66.4) && ((doy < 83) || (doy > 263))) {
			newjd = findNextSunset(JD, latitude, longitude)
			newtime = calcSunsetUTC(newjd, latitude, longitude)
			-(60 * zone) + daySavings
			if newtime > 1440 {
				newtime -= 1440
				newjd += 1.0
			}
			if newtime < 0 {
				newtime += 1440
				newjd -= 1.0
			}
			riseSetForm["sunset"].value =
				timeStringAMPMDate(newtime, newjd)
			riseSetForm["utcsunset"].value = "next sunset"
			riseSetForm["utcsolnoon"].value = ""

			// if Northern hemisphere and fall or winter, OR
			// if Southern hemisphere and spring or summer, use
			// next sunrise and last sunset

		} else if ((latitude > 66.4) && ((doy < 83) || (doy > 263))) || ((latitude < -66.4) && (doy > 79) && (doy < 267)) {
			newjd = findRecentSunset(JD, latitude, longitude)
			newtime = calcSunsetUTC(newjd, latitude, longitude)
			-(60 * zone) + daySavings
			if newtime > 1440 {
				newtime -= 1440
				newjd += 1.0
			}
			if newtime < 0 {
				newtime += 1440
				newjd -= 1.0
			}
			riseSetForm["sunset"].value =
				timeStringAMPMDate(newtime, newjd)
			riseSetForm["utcsunset"].value = "prior sunset"
			riseSetForm["solnoon"].value = "N/A"
			riseSetForm["utcsolnoon"].value = ""
		} else {
			alert("Cannot Find Sunset!")
		}
	}

}
*/
