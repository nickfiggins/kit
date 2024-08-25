package secrets

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Fetcher func(id string) (string, error)

// NewFetcher creates a new Fetcher that fetches secrets from AWS Secrets Manager
// or from environment variables. If the secret is not found in the environment,
// it will attempt to fetch it from Secrets Manager.
//
// It first looks up the environment variable with the key, if it's not found, it
// looks up the environment variable with the key prefixed with "SM_" to detect the token
// path.
func NewFetcher(ctx context.Context, getenv func(string) string) (Fetcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	sm := secretsmanager.NewFromConfig(cfg)

	return func(key string) (string, error) {
		if v := getenv(key); v != "" {
			return v, nil
		}

		smPath := fmt.Sprintf("SM_%s", key)
		v := getenv(fmt.Sprintf("SM_%s", key))
		if v == "" {
			return "", fmt.Errorf("can't fetch secret, both %q and %q not found in environment", key, smPath)
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		out, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: &v,
		})
		if err != nil {
			return "", fmt.Errorf("failed to get secret %q: %w", key, err)
		}
		return *out.SecretString, nil
	}, nil
}
