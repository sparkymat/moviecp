package tmdb

import (
	"context"
	"net/http"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

const (
	baseURL      = "https://api.themoviedb.org/3"
	imageBaseURL = "https://image.tmdb.org/t/p/w500"
)

func New(httpClient HTTPClient, redisClient RedisClient, apiToken string) *Provider {
	return &Provider{
		httpClient:  httpClient,
		apiToken:    apiToken,
		redisClient: redisClient,
	}
}

type Provider struct {
	httpClient  HTTPClient
	apiToken    string
	redisClient RedisClient
}
