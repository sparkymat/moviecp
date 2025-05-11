//nolint:tagliatelle
package tmdb

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/samber/lo"
	"github.com/sparkymat/moviecp/internal/provider"
)

type searchTVResponseItem struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	Overview         string  `json:"overview"`
	GenreIDs         []int64 `json:"genre_ids"`
	PosterPath       string  `json:"poster_path"`
	VoteAverage      float64 `json:"vote_average"`
	OriginalLanguage string  `json:"original_language"`
}

type searchTVResponse struct {
	Results    []searchTVResponseItem `json:"results"`
	TotalPages int64                  `json:"total_pages"`
}

func (p *Provider) SearchTV(ctx context.Context, query string, page int64) ([]TvShow, int64, error) {
	searchURL := fmt.Sprintf("%s/search/tv?api_key=%s&query=%s&page=%d", baseURL, p.apiToken, url.QueryEscape(query), page)

	response, err := provider.FetchJSON[searchTVResponse](ctx, p.redisClient, p.httpClient, searchURL)
	if err != nil {
		return []TvShow{}, 0, fmt.Errorf("failed to fetch tv shows: %w", err)
	}

	genresMap, err := p.ListGenres(ctx, "en") // Get genres names in en
	if err != nil {
		return []TvShow{}, 0, fmt.Errorf("failed to fetch genres: %w", err)
	}

	tvShows := lo.Map(response.Results, func(item searchTVResponseItem, _ int) TvShow {
		genres := lo.Map(item.GenreIDs, func(genreID int64, _ int) string { return genresMap[genreID] })

		return TvShow{
			ID:          strconv.FormatInt(item.ID, 10),
			Title:       item.Name,
			Genres:      genres,
			Overview:    item.Overview,
			PosterPath:  fmt.Sprintf("%s%s", imageBaseURL, item.PosterPath),
			Language:    item.OriginalLanguage,
			VoteAverage: item.VoteAverage,
		}
	})

	return tvShows, response.TotalPages, nil
}
