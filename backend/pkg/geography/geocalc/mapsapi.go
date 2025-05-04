package geocalc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

var mapsApiKey string

func init() {
	mapsApiKey := os.Getenv("MAPS_API_KEY")
	if mapsApiKey == "" {
		logrus.Fatal("MAPS_API_KEY environment variable is not set")
	}
}

func distanceMatrix(guess, answer string) float64 {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s&destinations=%s&key=%s",
		guess, answer, mapsApiKey)

	results, err := mapsApi(url)
	if err != nil {
		logrus.Fatalf("Failed to get distance matrix from Google Maps API: %v", err)
	}

	rows := results["rows"].([]interface{})[0].(map[string]interface{})
	elements := rows["elements"].([]interface{})[0].(map[string]interface{})
	distance := elements["distance"].(map[string]interface{})
	value := distance["value"].(float64) //distance in meters

	distanceKM := value / 1000
	return distanceKM
}

func geocode(country string) (float64, float64) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s",
		country, mapsApiKey)

	results, err := mapsApi(url)
	if err != nil {
		logrus.Fatalf("Failed to get geocode from Google Maps API: %v", err)
	}

	result := results["results"].([]interface{})[0].(map[string]interface{})
	geometry := result["geometry"].(map[string]interface{})
	location := geometry["location"].(map[string]interface{})
	lat := location["lat"].(float64)
	lng := location["lng"].(float64)

	return lat, lng
}

func mapsApi(url string) (map[string]interface{}, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Fatal("Error creating request to Google Maps API")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logrus.Fatalf("Failed to get a response from Google Maps API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logrus.Fatalf("Failed to make request to Google Maps API: %v", err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Fatalf("Failed to read response body from Google Maps API: %v", err)
	}

	var results map[string]interface{}
	if err := json.Unmarshal(data, &results); err != nil {
		logrus.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	return results, nil
}
