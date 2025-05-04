package geocalc

import "fmt"

func GetDistanceAndDirection(guess string, answer string) (float64, string, error) {
	guessLat, guessLng, err := geocode(guess)
	if err != nil {
		return 0, "", fmt.Errorf("failed to geocode guess: %w", err)
	}

	answerLat, answerLng, err := geocode(answer)
	if err != nil {
		return 0, "", fmt.Errorf("failed to geocode answer: %w", err)
	}

	distance, err := haversine(guessLat, guessLng, answerLat, answerLng)
	if err != nil {
		return 0, "", fmt.Errorf("failed to calculate haversine distance: %w", err)
	}

	distanceKM := distance / 1000
	direction := getCompass(guessLat, guessLng, answerLat, answerLng)

	return distanceKM, direction, nil
}
