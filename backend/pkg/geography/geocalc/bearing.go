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
		"N":   {354.375, 5.625}, // Narrow range for North
		"NNE": {5.625, 28.125},
		"NE":  {28.125, 61.875},
		"ENE": {61.875, 84.375},
		"E":   {84.375, 95.625}, // Narrow range for East
		"ESE": {95.625, 118.125},
		"SE":  {118.125, 151.875},
		"SSE": {151.875, 174.375},
		"S":   {174.375, 185.625}, // Narrow range for South
		"SSW": {185.625, 208.125},
		"SW":  {208.125, 241.875},
		"WSW": {241.875, 264.375},
		"W":   {264.375, 275.625}, // Narrow range for West
		"WNW": {275.625, 298.125},
		"NW":  {298.125, 331.875},
		"NNW": {331.875, 354.375},
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
	guessLatRad := guessLat * (math.Pi / 180)
	guessLngRad := guessLng * (math.Pi / 180)
	answerLatRad := answerLat * (math.Pi / 180)
	answerLngRad := answerLng * (math.Pi / 180)
	deltaLng := answerLngRad - guessLngRad

	y := math.Sin(deltaLng) * math.Cos(answerLatRad)
	x := math.Cos(guessLatRad)*math.Sin(answerLatRad) - math.Sin(guessLatRad)*math.Cos(answerLatRad)*math.Cos(answerLngRad-guessLngRad)
	bearing := math.Atan2(y, x)

	return bearing
}
