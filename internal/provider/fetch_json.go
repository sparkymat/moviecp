package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

const (
	UserAgent          = "tube/1.0"
	HTTPTimeoutSeconds = 20
)

var ErrRequestFailed = errors.New("request failed")

type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

func FetchJSON[T any](ctx context.Context, redisClient RedisClient, httpClient HTTPClient, url string) (T, error) {
	log.Infof("fetching JSON from %s", url)

	var response T

	if redisClient != nil {
		cachedData, err := redisClient.Get(ctx, url).Result()
		if err == nil {
			if err = json.Unmarshal([]byte(cachedData), &response); err != nil {
				return response, fmt.Errorf("failed to unmarshal cached response. err: %w", err)
			}

			log.Info("returning cached data")

			return response, nil
		}
	}

	ctx, cancel := context.WithTimeout(ctx, HTTPTimeoutSeconds*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response, fmt.Errorf("unable to form request. err: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("request failed. err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("failed to fetch response. err: %w", ErrRequestFailed)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response. err: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("request failed with %d. err: %w", resp.StatusCode, ErrRequestFailed)
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return response, fmt.Errorf("failed to unmarshal response. err: %w", err)
	}

	if redisClient != nil {
		if err := redisClient.Set(ctx, url, string(respBody), 24*time.Hour).Err(); err != nil { //nolint:mnd
			return response, fmt.Errorf("failed to cache response. err: %w", err)
		}
	}

	return response, nil
}
