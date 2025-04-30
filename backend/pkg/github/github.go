package github

import (
	"fmt"
	"os"
)

type GitConfig struct {
	Token  string `env:"GITHUB_TOKEN"`
	Repo   string `env:"GITHUB_REPO"`
	Owner  string `env:"GITHUB_OWNER"`
	Branch string `env:"GITHUB_BRANCH"`
}

func LoadGitConfig() (*GitConfig, error) {
	gitConfig := &GitConfig{
		Token:  os.Getenv("GITHUB_TOKEN"),
		Repo:   os.Getenv("GITHUB_REPO"),
		Owner:  os.Getenv("GITHUB_OWNER"),
		Branch: os.Getenv("GITHUB_BRANCH"),
	}

	if gitConfig.Token == "" || gitConfig.Repo == "" || gitConfig.Owner == "" || gitConfig.Branch == "" {
		return nil, fmt.Errorf("missing required GitHub configuration")
	}

	return gitConfig, nil
}
