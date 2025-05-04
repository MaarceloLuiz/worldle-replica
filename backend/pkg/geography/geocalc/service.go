package geocalc

import "fmt"

func GetDistance(guess string, answer string) (float64, error) {
	distance, err := distanceMatrix(guess, answer)
	if err != nil {
		return 0, fmt.Errorf("failed to get distance from Google Maps API: %w", err)
	}

	distanceKM := distance / 1000
	return distanceKM, nil
}

func GetDirection(guess string, answer string) (string, error) {
	guessLat, guessLng, err := geocode(guess)
	if err != nil {
		return "", fmt.Errorf("failed to geocode guess: %w", err)
	}

	answerLat, answerLng, err := geocode(answer)
	if err != nil {
		return "", fmt.Errorf("failed to geocode answer: %w", err)
	}

	return getCompass(guessLat, guessLng, answerLat, answerLng), nil
}
