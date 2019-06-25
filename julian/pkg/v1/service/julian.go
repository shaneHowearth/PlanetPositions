package julian

import (
	"fmt"
	"math"
)

/*----------------------------------------------------------------------
**
**  Copyright (C) 2018
**  Standards Of Fundamental Astronomy Board
**  of the International Astronomical Union.
**
**  =====================
**  SOFA Software License
**  =====================
**
**  NOTICE TO USER:
**
**  BY USING THIS SOFTWARE YOU ACCEPT THE FOLLOWING SIX TERMS AND
**  CONDITIONS WHICH APPLY TO ITS USE.
**
**  1. The Software is owned by the IAU SOFA Board ("SOFA").
**
**  2. Permission is granted to anyone to use the SOFA software for any
**     purpose, including commercial applications, free of charge and
**     without payment of royalties, subject to the conditions and
**     restrictions listed below.
**
**  3. You (the user) may copy and distribute SOFA source code to others,
**     and use and adapt its code and algorithms in your own software,
**     on a world-wide, royalty-free basis.  That portion of your
**     distribution that does not consist of intact and unchanged copies
**     of SOFA source code files is a "derived work" that must comply
**     with the following requirements:
**
**     a) Your work shall be marked or carry a statement that it
**        (i) uses routines and computations derived by you from
**        software provided by SOFA under license to you; and
**        (ii) does not itself constitute software provided by and/or
**        endorsed by SOFA.
**
**     b) The source code of your derived work must contain descriptions
**        of how the derived work is based upon, contains and/or differs
**        from the original SOFA software.
**
**     c) The names of all routines in your derived work shall not
**        include the prefix "iau" or "sofa" or trivial modifications
**        thereof such as changes of case.
**
**     d) The origin of the SOFA components of your derived work must
**        not be misrepresented;  you must not claim that you wrote the
**        original software, nor file a patent application for SOFA
**        software or algorithms embedded in the SOFA software.
**
**     e) These requirements must be reproduced intact in any source
**        distribution and shall apply to anyone to whom you have
**        granted a further right to modify the source code of your
**        derived work.
**
**     Note that, as originally distributed, the SOFA software is
**     intended to be a definitive implementation of the IAU standards,
**     and consequently third-party modifications are discouraged.  All
**     variations, no matter how minor, must be explicitly marked as
**     such, as explained above.
**
**  4. You shall not cause the SOFA software to be brought into
**     disrepute, either by misuse, or use for inappropriate tasks, or
**     by inappropriate modification.
**
**  5. The SOFA software is provided "as is" and SOFA makes no warranty
**     as to its use or performance.   SOFA does not and cannot warrant
**     the performance or results which the user may obtain by using the
**     SOFA software.  SOFA makes no warranties, express or implied, as
**     to non-infringement of third party rights, merchantability, or
**     fitness for any particular purpose.  In no event will SOFA be
**     liable to the user for any consequential, incidental, or special
**     damages, including any lost profits or lost savings, even if a
**     SOFA representative has been advised of such damages, or for any
**     claim by any third party.
**
**  6. The provision of any version of the SOFA software under the terms
**     and conditions specified herein does not imply that future
**     versions will also be made available under the same terms and
**     conditions.
*
**  In any published work or commercial product which uses the SOFA
**  software directly, acknowledgement (see www.iausofa.org) is
**  appreciated.
**
**  Correspondence concerning SOFA software should be addressed as
**  follows:
**
**      By email:  sofa@ukho.gov.uk
**      By post:   IAU SOFA Center
**                 HM Nautical Almanac Office
**                 UK Hydrographic Office
**                 Admiralty Way, Taunton
**                 Somerset, TA1 2DN
**                 United Kingdom
**
**--------------------------------------------------------------------*/
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
	// This algorithm is valid from -4800 March, but rejects dates before -4799 January 1
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

	monthYear := float64(month-14) / 12
	// Wikipedia
	//JDN = (1461 × (Y + 4800 + (M − 14)/12))/4 +(367 × (M − 2 − 12 × ((M − 14)/12)))/12 − (3 × ((Y + 4900 + (M - 14)/12)/100))/4 + D − 32075
	julianDay := (1461 * (float64(year) + 4800 + monthYear)) / 4
	fmt.Println(julianDay)
	julianDay += (367 * (float64(month) - 2 - 12*monthYear)) / 12
	fmt.Println(julianDay)
	julianDay -= (3*((float64(year+4900)+monthYear)/100))/4 + float64(day-32075)
	//julianDay -= 32054
	return int32(julianDay), nil
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
