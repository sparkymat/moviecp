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

type searchMoviesResponseItem struct {
	ID               int64   `json:"id"`
	Title            string  `json:"title"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	GenreIDs         []int64 `json:"genre_ids"`
	PosterPath       string  `json:"poster_path"`
	VoteAverage      float64 `json:"vote_average"`
	OriginalLanguage string  `json:"original_language"`
}

type searchMoviesResponse struct {
	Results    []searchMoviesResponseItem `json:"results"`
	TotalPages int64                      `json:"total_pages"`
}

func (p *Provider) SearchMovies(ctx context.Context, query string, page int64) ([]Movie, int64, error) {
	searchURL := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s&page=%d", baseURL, p.apiToken, url.QueryEscape(query), page)

	response, err := provider.FetchJSON[searchMoviesResponse](ctx, p.redisClient, p.httpClient, searchURL)
	if err != nil {
		return []Movie{}, 0, fmt.Errorf("failed to fetch movies: %w", err)
	}

	genresMap, err := p.ListGenres(ctx, "en") // Get genres names in en
	if err != nil {
		return []Movie{}, 0, fmt.Errorf("failed to fetch genres: %w", err)
	}

	movies := lo.Map(response.Results, func(item searchMoviesResponseItem, _ int) Movie {
		genres := lo.Map(item.GenreIDs, func(genreID int64, _ int) string { return genresMap[genreID] })

		return Movie{
			ID:          strconv.FormatInt(item.ID, 10),
			Title:       item.Title,
			Genres:      genres,
			Overview:    item.Overview,
			PosterPath:  fmt.Sprintf("%s%s", imageBaseURL, item.PosterPath),
			ReleaseDate: item.ReleaseDate,
			Language:    item.OriginalLanguage,
			VoteAverage: item.VoteAverage,
		}
	})

	return movies, response.TotalPages, nil
}
