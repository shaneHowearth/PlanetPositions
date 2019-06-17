package julian

import "math"

const jan12000 = float64(2451545)
const century = float64(36525)

// GetJulianDay -
func GetJulianDay(year, month, day int32, universalTime float64) float64 {
	// This function converts a gregorian date/time to a Julian day
	// universal_time is comprised of hours and fraction of hours
	// eg. 7:30am would be 7.5
	// Source: http://aa.usno.navy.mil/faq/docs/JD_Formula.php
	// This formula is only good for the years between 1801–2099
	//
	fday := float64(day)
	fmonth := float64(month)
	fyear := float64(year)
	var sign float64
	if (100.0*fyear + fmonth - 190002.5) > 0 {
		sign = 1.0
	} else {
		sign = -1.0
	}

	// we only want the integer part of these results
	//f1 := float64(int(7 * (fyear + float64(int((fmonth+9)/12))) / 4))
	//f2 := float64(int((275 * fmonth) / 9))

	//julian_day := 367.0*fyear - f1 + f2 + fday + 1721013.5 + universal_time/24.0 - 0.5*sign + 0.5

	f0 := fyear + 4800.0 + (fmonth-14.0)/12.0
	f1 := (1461.0 * f0) / 4.0
	f2a := (fmonth - 14.0) / 12.0
	f2b := (367.0 * (fmonth - 2.0 - 12.0*f2a)) / 12.0
	f3a := fyear + 4900.0 + (fmonth-14.0)/12.0
	f3b := 3.0 * (f3a / 100.0)
	f3c := f3b / 4.0
	julianDay := float64(int(f1+f2b-f3c+fday-32075.0)) + universalTime/24.0 + sign*0.5
	return julianDay
}

// CalcDayOfWeek -
func CalcDayOfWeek(julianDay float64) int32 {
	// Sunday is 0
	a := int32((julianDay + float64(1.5))) % int32(7)
	return a
}

// DayFromJulianDay -
func DayFromJulianDay(julianDay float64) (year, month, day int32) {
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

// TimeJulianCentury -
func TimeJulianCentury(julianDay float64) float64 {
	return (julianDay - jan12000) / century
}

// JulianDayFromJulianCentury -
func JulianDayFromJulianCentury(t float64) float64 {
	return t*century + jan12000
}
