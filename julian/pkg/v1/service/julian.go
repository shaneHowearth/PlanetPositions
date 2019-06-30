package julian

import (
	"fmt"
	"math"
)

const jan12000 = float64(2451545)
const century = float64(36525)

func isLeapYear(year int32) bool {
	if year%4 == 0 {
		if year%100 == 0 {
			return year%400 == 0
		}
		return true
	}
	return false
}

// GetJulianDay -
func GetJulianDay(year, month, day int32, universalTime float64) (int32, error) {
	// Not going to support BC values until I can figure out where the bug in the algorithm is
	if year < 1 {
		return 0, fmt.Errorf("BC dates not currently supported")
	}
	if year < -4712 {
		return int32(0), fmt.Errorf("dates earlier than January 1 4799 BC will not be computed")
	}
	// Validate month
	if month < 1 || month > 12 {
		return int32(0), fmt.Errorf("received an impossible month number %d", month)
	}

	// Validate day
	monthDays := [12]int32{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	maxDay := monthDays[month-1]

	//leapYear := 0
	if month == 2 {
		if isLeapYear(year) {
			//leapYear = 1
			maxDay = 29
		}
	}
	if day > maxDay {
		return int32(0), fmt.Errorf("received an impossible day number: %d, max possible days for that month is: %d", day, maxDay)
	}

	//monthYear := int32(float64(month-14) / 12)
	// Wikipedia
	//JDN = (1461 × (Y + 4800 + (M − 14)/12))/4 +(367 × (M − 2 − 12 × ((M − 14)/12)))/12 − (3 × ((Y + 4900 + (M - 14)/12)/100))/4 + D − 32075
	julianDay := (1461 * (year + 4800 + (month-14)/12)) / 4
	fmt.Println(julianDay)
	julianDay += (367 * (month - 2 - 12*((month-14)/12))) / 12
	fmt.Println(julianDay)
	julianDay -= (3*((year+4900+(month-14)/12)/100))/4 + day - 32075
	//julianDay -= 32054
	julianDay -= 64107
	return julianDay, nil
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
	C := math.Floor((B - 122.1) / yearLen)
	D := math.Floor(yearLen * C)
	E := math.Floor((B - D) / monthLen)

	day = int32(B - D - math.Floor(monthLen*E) + f)
	month = int32(E) - int32(1)
	if !(E < float64(14)) {
		month -= int32(12)
	}
	year = int32(C - 4715)
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
func iJulianDayFromJulianCentury(t float64) float64 {
	return t*century + jan12000
}
