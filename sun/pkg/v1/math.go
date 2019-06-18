package v1

import "math"

func radiansToDegrees(angleRad float64) float64 {
	return 180 * angleRad / math.Pi
}

func degreesToRadians(angleDeg float64) float64 {
	return math.Pi * angleDeg / 180.0
}
