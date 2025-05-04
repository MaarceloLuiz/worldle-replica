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
	mapsApiKey = os.Getenv("MAPS_API_KEY")
	if mapsApiKey == "" {
		logrus.Fatal("MAPS_API_KEY environment variable is not set")
	}
}

func geocode(country string) (float64, float64, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s",
		country, mapsApiKey)

	results, err := mapsApi(url)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get geocode from Google Maps API: %w", err)
	}

	result := results["results"].([]interface{})[0].(map[string]interface{})
	geometry := result["geometry"].(map[string]interface{})
	location := geometry["location"].(map[string]interface{})
	lat := location["lat"].(float64)
	lng := location["lng"].(float64)

	return lat, lng, nil
}

func mapsApi(url string) (map[string]interface{}, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to Google Maps API: %w", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get a response from Google Maps API: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google Maps API returned status: %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from Google Maps API: %w", err)
	}

	var results map[string]interface{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return results, nil
}
