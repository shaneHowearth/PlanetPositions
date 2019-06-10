package julian

// GetJulianDay -
func GetJulianDay(year, month, day int, universalTime float64) float64 {
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
