package geocalc

import (
	"math"

	"github.com/sirupsen/logrus"
)

func getCompass(guessLat, guessLng, answerLat, answerLng float64) string {
	bearing := initialBearing(guessLat, guessLng, answerLat, answerLng)
	degrees := radiansToDegrees(bearing)
	normalizedDegrees := normalizeDegrees(degrees)
	compass := degreesToCompass(normalizedDegrees)
	if compass == "" {
		logrus.Fatalf("Invalid compass direction for degrees: %f", normalizedDegrees)
	}

	return compass
}

func degreesToCompass(degrees float64) string {
	ranges := map[string][2]float64{
		"N":  {337.5, 22.5},
		"NE": {22.5, 67.5},
		"E":  {67.5, 112.5},
		"SE": {112.5, 157.5},
		"S":  {157.5, 202.5},
		"SW": {202.5, 247.5},
		"W":  {247.5, 292.5},
		"NW": {292.5, 337.5},
	}

	// special case for North since it wraps around 360
	if degrees > 337.5 || degrees <= 22.5 {
		return "N"
	}

	for direction, bounds := range ranges {
		if bounds[0] < degrees && degrees <= bounds[1] {
			return direction
		}
	}

	return "" // should never reach here if the input is valid
}

// When calculate the initial bearing using the atan2 function, the result (after converting to degrees)
// can be any value between -180° and +180°. However, compass bearings are always expressed as 0° to 360°.
// normalizeDegrees adjusts an angle to ensure it falls within the range of 0° to 360°.
func normalizeDegrees(degrees float64) float64 {
	return math.Mod(degrees+360, 360)
}

func radiansToDegrees(bearing float64) float64 {
	return bearing * (180.0 / math.Pi)
}

// θ = atan2( sin Δλ ⋅ cos φ2 , cos φ1 ⋅ sin φ2 − sin φ1 ⋅ cos φ2 ⋅ cos Δλ )
func initialBearing(guessLat, guessLng, answerLat, answerLng float64) float64 {
	y := math.Sin(answerLng-guessLng) * math.Cos(answerLat)
	x := math.Cos(guessLat)*math.Sin(answerLat) - math.Sin(guessLat)*math.Cos(answerLat)*math.Cos(answerLng-guessLng)
	bearing := math.Atan2(y, x)

	return bearing
}
