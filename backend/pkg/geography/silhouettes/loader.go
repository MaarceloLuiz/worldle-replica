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
	gitConfig, err := github.LoadGitConfig()
	if err != nil {
		logrus.Error("Error loading GitHub configuration")
		return nil, err
	}

	country, err := getRandomCountry()
	if err != nil {
		logrus.Error("Error generating random country")
		return nil, err
	}

	apiURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contents/silhouettes/%s.png?ref=%s",
		gitConfig.Owner, gitConfig.Repo, country, gitConfig.Branch,
	)

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		logrus.Error("Error creating request to GitHub API")
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+gitConfig.Token)
	request.Header.Set("Accept", "application/vnd.github.raw")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logrus.Error("Failed to make request to GitHub API: Authorization issue or network error")
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logrus.Error("Failed to fetch silhouette from GitHub API: ", response.Status)
		return nil, fmt.Errorf("failed to fetch silhouette from GitHub API: %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error("Failed to read response body from GitHub API")
		return nil, err
	}

	return data, nil
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
