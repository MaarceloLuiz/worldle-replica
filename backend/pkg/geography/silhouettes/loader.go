package silhouettes

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/MaarceloLuiz/worldle-replica/pkg/github"
	"github.com/sirupsen/logrus"
)

type SilhouetteResponse struct {
	URL string `json:"url"`
}

func FetchSilhouette() ([]byte, error) {
	gitConfig, err := github.LoadGitConfig()
	if err != nil {
		logrus.Error("Error loading GitHub configuration")
		return nil, err
	}

	country, err := getRandomCountry()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/silhouettes/%s.png",
		gitConfig.Owner, gitConfig.Repo, gitConfig.Branch, country)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("Error creating request to GitHub API")
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+gitConfig.Token)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logrus.Error("Failed to make request to GitHub API: Authorization issue or network error")
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logrus.Error("Failed to fetch silhouette from GitHub API: ", response.Status)
		return nil, fmt.Errorf("failed to fetch silhouette from GitHub API: %s", response.Status)
	}

	silhouetteResponse := SilhouetteResponse{URL: url}
	jsonResponse, err := json.Marshal(silhouetteResponse)
	if err != nil {
		logrus.Error("Failed to marshal silhouette response to JSON")
		return nil, err
	}

	return jsonResponse, nil
}

func getRandomCountry() (string, error) {
	gitConfig, err := github.LoadGitConfig()
	if err != nil {
		logrus.Error("Error loading GitHub configuration")
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/silhouettes?ref=%s",
		gitConfig.Owner, gitConfig.Repo, gitConfig.Branch)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("Error creating request to GitHub API")
		return "", err
	}

	request.Header.Set("Authorization", "Bearer "+gitConfig.Token)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logrus.Error("Failed to make request to GitHub API: Authorization issue or network error")
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logrus.Error("Failed to fetch silhouettes from GitHub API: ", response.Status)
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
		logrus.Error("No silhouette files found in the repository")
		return "", fmt.Errorf("no silhouette files found in the repository")
	}

	seed := time.Now().UnixNano() // to avoid seeding the same number every time
	random := rand.New(rand.NewSource(seed))

	randomIndex := random.Intn(len(countries))
	randomCountry := countries[randomIndex]

	return randomCountry, nil
}
