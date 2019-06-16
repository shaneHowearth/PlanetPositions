package v1

import "math"

func radiansToDegrees(angleRad float64) float64 {
	return 180 * angleRad / math.Pi
}

func degreesToRadians(angleDeg float64) float64 {
	return math.Pi * angleDeg / 180.0
}

func calcDayOfYear(month, day int32, leapYear bool) int32 {
	k := int32(1)
	if !leapYear {
		k = 2
	}
	dayOfYear := int32(math.Floor(float64(275*month))/float64(9)) - k*int32(math.Floor(float64(month+9)/float64(12))) + day - int32(30)
	return dayOfYear
}

func calcDayOfWeek(julianDay float64) int32 {
	// Sunday is 0
	a := int32((julianDay + float64(1.5))) % int32(7)
	return a
}

func dayFromJulianDay(julianDay float64) (year, month, day int32) {
	z := math.Floor(julianDay + float64(0.5))
	f := (julianDay + float64(0.5)) - z
	A := z
	if !(z < float64(2299161)) {
		alpha := math.Floor((z - float64(1867216.25)) / float64(36524.25))
		A = z + float64(1) + alpha - math.Floor(alpha/float64(4))
	}
	yearLen := float64(365.25)
	monthLen := float64(30.6001)
	B := A + float64(1524)
	C := math.Floor((B - float64(122.1)) / yearLen)
	D := math.Floor(yearLen * C)
	E := math.Floor((B - D) / monthLen)

	day = int32(B - D - math.Floor(monthLen*E) + f)
	month = int32(E) - int32(1)
	if !(E < float64(14)) {
		month -= int32(12)
	}
	year = int32(C) - int32(4715)
	if month > 2 {
		year = year - 1
	}
	return year, month, day
}

func calcGeomMeanLongSun(numCenturies float64) {

}

const jan12000 = float64(2451545)
const century = float64(36525)

func timeJulianCentury(julianDay float64) float64 {
	return (julianDay - jan12000) / century
}

func julianDayFromJulianCentury(t float64) float64 {
	return t*century + jan12000
}

func geometricMeanLongitudeSun(t float64) float64 {
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

func geometricMeanAnamolySun(t float64) float64 {
	// t is the number of Julian centuries since J2000.0
	return 357.52911 + t*(35999.05029-0.0001537*t)
}

func eccentricityEarthOrbit(t float64) float64 {
	return 0.016708634 - t*(0.000042037+0.0000001267*t)
}

func equationCentreSun(t float64) float64 {
	m := geometricMeanAnamolySun(t)
	mrad := degreesToRadians(m)
	sinm := math.Sin(mrad)
	sin2m := math.Sin(mrad + mrad)
	sin3m := math.Sin(mrad + mrad + mrad)

	return sinm*(1.914602-t*(0.004817+0.000014*t)) + sin2m*(0.019993-0.000101*t) + sin3m*0.000289 // In Degrees
}

func sunTrueLongitude(t float64) float64 {
	l0 := geometricMeanLongitudeSun(t)
	c := equationCentreSun(t)

	return l0 + c // In Degrees
}

func sunTrueAnamoly(t float64) float64 {
	m := geometricMeanAnamolySun(t)
	c := equationCentreSun(t)

	return m + c // In Degrees

}

func sunRadiusVector(t float64) float64 {
	v := sunTrueAnamoly(t)
	e := eccentricityEarthOrbit(t)

	return (1.000001018 * (1 - e*e)) / (1 + e*math.Cos(degreesToRadians(v))) // In Astronomical Units
}

func sunApparentLongitude(t float64) float64 {
	o := sunTrueLongitude(t)

	omega := 125.04 - 1934.136*t
	return o - 0.00569 - 0.00478*math.Sin(degreesToRadians(omega)) // In Degrees
}

func meanObliquityOfEcliptic(t float64) float64 {

	seconds := 21.448 - t*(46.8150+t*(0.00059-t*(0.001813)))
	return 23.0 + (26.0+(seconds/60.0))/60.0 // In Degrees
}

func obliquityCorrection(t float64) float64 {

	e0 := meanObliquityOfEcliptic(t)

	omega := 125.04 - 1934.136*t
	return e0 + 0.00256*math.Cos(degreesToRadians(omega))
}

func sunRtAscension(t float64) float64 {
	e := obliquityCorrection(t)
	lambda := sunApparentLongitude(t)

	tananum := (math.Cos(degreesToRadians(e)) * math.Sin(degreesToRadians(lambda)))
	tanadenom := (math.Cos(degreesToRadians(lambda)))
	//TODO TEST THIS BECAUSE Atan2 might not be the same
	return radiansToDegrees(math.Atan2(tananum, tanadenom)) // In Degrees
}
func sunDeclination(t float64) float64 {
	e := obliquityCorrection(t)
	lambda := sunApparentLongitude(t)

	sint := math.Sin(degreesToRadians(e)) * math.Sin(degreesToRadians(lambda))
	return radiansToDegrees(math.Asin(sint)) // In Degrees
}
func equationOfTime(t float64) float64 {
	epsilon := obliquityCorrection(t)
	l0 := geometricMeanLongitudeSun(t)
	e := eccentricityEarthOrbit(t)
	m := geometricMeanAnamolySun(t)

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

func hourAngleSunrise(lat, solarDec float64) float64 {
	latRad := degreesToRadians(lat)
	sdRad := degreesToRadians(solarDec)

	// HAarg := (math.Cos(degreesToRadians(90.833))/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad))

	HA := (math.Acos(math.Cos(degreesToRadians(90.833))/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad)))

	return HA // In Radians
}

func hourAngleSunset(lat, solarDec float64) float64 {
	latRad := degreesToRadians(lat)
	sdRad := degreesToRadians(solarDec)

	// HAarg := (math.Cos(degreesToRadians(90.833))/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad))

	HA := (math.Acos(math.Cos(degreesToRadians(90.833))/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad)))

	return -HA // in radians
}
func solNoonUTC(t, longitude float64) float64 {
	// First pass uses approximate solar noon to calculate eqtime
	tnoon := timeJulianCentury(julianDayFromJulianCentury(t) + longitude/360.0)
	eqTime := equationOfTime(tnoon)
	solNoonUTC := 720 + (longitude * 4) - eqTime // min

	newt := timeJulianCentury(julianDayFromJulianCentury(t) - 0.5 + solNoonUTC/1440.0)

	eqTime = equationOfTime(newt)
	// var solarNoonDec = calcSunDeclination(newt)
	solNoonUTC = 720 + (longitude * 4) - eqTime // min

	return solNoonUTC
}

func sunriseUTC(JD, latitude, longitude float64) float64 {
	t := timeJulianCentury(JD)

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin := solNoonUTC(t, longitude)
	tnoon := timeJulianCentury(JD + noonmin/1440.0)

	// *** First pass to approximate sunrise (using solar noon)

	eqTime := equationOfTime(tnoon)
	solarDec := sunDeclination(tnoon)
	hourAngle := hourAngleSunrise(latitude, solarDec)

	delta := longitude - radiansToDegrees(hourAngle)
	timeDiff := 4 * delta              // in minutes of time
	timeUTC := 720 + timeDiff - eqTime // in minutes

	// alert("eqTime = " + eqTime + "\nsolarDec = " + solarDec + "\ntimeUTC = " + timeUTC)

	// *** Second pass includes fractional jday in gamma calc

	newt := timeJulianCentury(julianDayFromJulianCentury(t) + timeUTC/1440.0)
	eqTime = equationOfTime(newt)
	solarDec = sunDeclination(newt)
	hourAngle = hourAngleSunrise(latitude, solarDec)
	delta = longitude - radiansToDegrees(hourAngle)
	timeDiff = 4 * delta
	timeUTC = 720 + timeDiff - eqTime // in minutes

	// alert("eqTime = " + eqTime + "\nsolarDec = " + solarDec + "\ntimeUTC = " + timeUTC)

	return timeUTC
}

func sunsetUTC(JD, latitude, longitude float64) float64 {
	t := timeJulianCentury(JD)

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin := solNoonUTC(t, longitude)
	tnoon := timeJulianCentury(JD + noonmin/1440.0)

	// First calculates sunrise and approx length of day

	eqTime := equationOfTime(tnoon)
	solarDec := sunDeclination(tnoon)
	hourAngle := hourAngleSunset(latitude, solarDec)

	delta := longitude - radiansToDegrees(hourAngle)
	timeDiff := 4 * delta
	timeUTC := 720 + timeDiff - eqTime

	// first pass used to include fractional day in gamma calc

	newt := timeJulianCentury(julianDayFromJulianCentury(t) + timeUTC/1440.0)
	eqTime = equationOfTime(newt)
	solarDec = sunDeclination(newt)
	hourAngle = hourAngleSunset(latitude, solarDec)

	delta = longitude - radiansToDegrees(hourAngle)
	timeDiff = 4 * delta
	timeUTC = 720 + timeDiff - eqTime // in minutes

	return timeUTC
}
