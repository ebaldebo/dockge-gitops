package config

import "github.com/ebaldebo/dockge-gitops/internal/env"

type Config struct {
	RepoUrl     string
	Pat         string
	PollingRate string
}

func New() (*Config, error) {
	repoUrl, err := env.GetEnvVar(true, "REPO_URL", "")
	if err != nil {
		return nil, err
	}

	pat, err := env.GetEnvVar(false, "PAT", "")
	if err != nil {
		return nil, err
	}

	pollingRate, err := env.GetEnvVar(false, "POLLING_RATE", "5m")
	if err != nil {
		return nil, err
	}

	return &Config{
		RepoUrl:     repoUrl,
		Pat:         pat,
		PollingRate: pollingRate,
	}, err
}
