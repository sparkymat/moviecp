package tmdb

import (
	"context"
	"fmt"

	"github.com/sparkymat/moviecp/internal/provider"
)

type listGenresResponseItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type listGenresResponse struct {
	Genres []listGenresResponseItem `json:"genres"`
}

func (p *Provider) ListGenres(ctx context.Context, language string) (map[int64]string, error) {
	url := fmt.Sprintf("%s/genre/movie/list?api_key=%s&language=%s", baseURL, p.apiToken, language)

	response, err := provider.FetchJSON[listGenresResponse](ctx, p.redisClient, p.httpClient, url)
	if err != nil {
		return map[int64]string{}, fmt.Errorf("failed to fetch genres: %w", err)
	}

	genresMap := map[int64]string{}

	for _, item := range response.Genres {
		genresMap[item.ID] = item.Name
	}

	return genresMap, nil
}
