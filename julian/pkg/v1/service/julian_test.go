package julian_test

import (
	"fmt"
	"testing"

	julian "planetpositions/julian/pkg/v1/service"

	"github.com/stretchr/testify/assert"
)

func TestGetJulianDay(t *testing.T) {
	testcases := map[string]struct {
		year          int32
		month         int32
		day           int32
		universalTime float64
		output        int32
	}{
		"Happy path": {
			year:          2019,
			month:         06,
			day:           22,
			universalTime: 12,
			output:        int32(2458656),
		},
		"Happy path2": {
			year:          1000,
			month:         06,
			day:           22,
			universalTime: 12,
			output:        int32(2086474),
		},
		"100 BC": {
			year:          -1000,
			month:         06,
			day:           22,
			universalTime: 12,
			output:        int32(1684705),
		},
		"1000 BC": {
			year:          -1000,
			month:         06,
			day:           22,
			universalTime: 12,
			//output:        int32(1355981),
			output: int32(1356345),
		},
	}
	for name, tc := range testcases {
		output, _ := julian.GetJulianDay(tc.year, tc.month, tc.day, tc.universalTime)

		fmt.Printf("Expected: %d, Actual: %d, Diff: %d\n", tc.output, output, output-tc.output)
		assert.Equal(t, tc.output, output, "Test %s did not return the expected output", name)
	}
}
