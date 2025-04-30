package silhouettes

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/MaarceloLuiz/worldle-replica/pkg/github"
	"github.com/sirupsen/logrus"
)

func FetchSilhouette() ([]byte, error) {
	country, err := getRandomCountry()
	if err != nil {
		logrus.Error("Error generating random country")
		return nil, err
	}

	endpoint := fmt.Sprintf("contents/silhouettes/%s.png", country)
	headers := map[string]string{
		"Accept": "application/vnd.github.raw",
	}

	request, err := createGitHubRequest(endpoint, headers)
	if err != nil {
		logrus.Error("Error creating new GitHub request")
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logrus.Errorf("Failed to make request to GitHub API - Authorization issue or network error: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch silhouette from GitHub API")
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("Failed to read response body from GitHub API: %v", err)
		return nil, err
	}

	return data, nil
}

func getRandomCountry() (string, error) {
	request, err := createGitHubRequest("contents/silhouettes", nil)
	if err != nil {
		logrus.Error("Error creating new GitHub request")
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logrus.Errorf("Failed to make request to GitHub API - Authorization issue or network error: %v", err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch silhouettes from GitHub API: %s", response.Status)
	}

	var files []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(response.Body).Decode(&files); err != nil {
		logrus.Error("Failed to decode response from GitHub API")
		return "", err
	}

	var countries []string
	for _, file := range files {
		if strings.HasSuffix(file.Name, ".png") {
			countries = append(countries, strings.TrimSuffix(file.Name, ".png"))
		}
	}

	if len(countries) == 0 {
		return "", fmt.Errorf("no silhouette files found in the repository")
	}

	seed := time.Now().UnixNano() // to avoid seeding the same number every time
	random := rand.New(rand.NewSource(seed))

	randomIndex := random.Intn(len(countries))
	randomCountry := countries[randomIndex]

	return randomCountry, nil
}

func createGitHubRequest(endpoint string, headers map[string]string) (*http.Request, error) {
	gitConfig, err := github.LoadGitConfig()
	if err != nil {
		logrus.Error("Error loading GitHub configuration")
		return nil, err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/%s?ref=%s",
		gitConfig.Owner, gitConfig.Repo, endpoint, gitConfig.Branch)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("Error creating request to GitHub API")
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+gitConfig.Token)
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	return request, nil
}
