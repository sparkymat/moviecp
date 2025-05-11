package movie

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"github.com/sparkymat/moviecp/internal/provider/tmdb"
)

type SearchMoviesResultItem struct {
	Title       string   `json:"title"`
	ReleaseDate string   `json:"releaseData"`
	Synopsis    string   `json:"synopsis"`
	Genres      []string `json:"genres"`
	Language    string   `json:"language"`
}

type SearchMoviesResponse struct {
	Results []SearchMoviesResultItem
}

func (s *Service) SearchMovies(ctx context.Context, query string) ([]byte, error) {
	movies, _, err := s.tmdb.SearchMovies(ctx, query, 1)
	if err != nil {
		return []byte(`{
			"error": "failed to query"
}`), fmt.Errorf("failed to query movies: %w", err)
	}

	resultItems := lo.Map(movies, func(m tmdb.Movie, _ int) SearchMoviesResultItem {

		i := SearchMoviesResultItem{
			Title:       m.Title,
			ReleaseDate: m.ReleaseDate,
			Synopsis:    m.Overview,
			Genres:      m.Genres,
			Language:    m.Language,
		}

		return i
	})

	response := SearchMoviesResponse{
		Results: resultItems,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return []byte(`
		"error": "failed to marshal response"
`), fmt.Errorf("failed to marshal response: %w", err)
	}

	return responseBytes, nil
}
