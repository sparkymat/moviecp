package movie

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"github.com/sparkymat/moviecp/internal/provider/tmdb"
)

type SearchArtistsResultItem struct {
	ID         string                   `json:"id"`
	Department string                   `json:"department"`
	Name       string                   `json:"name"`
	Gender     string                   `json:"gender"`
	Movies     []SearchMoviesResultItem `json:"movies"`
}

type SearchArtistsResponse struct {
	Results []SearchArtistsResultItem
}

func (s *Service) SearchArtists(ctx context.Context, query string) ([]byte, error) {
	artists, _, err := s.tmdb.SearchArtists(ctx, query, 1)
	if err != nil {
		return []byte(`{"error": "failed to query"}`), fmt.Errorf("failed to query movies: %w", err)
	}

	resultItems := lo.Map(artists, func(m tmdb.Artist, _ int) SearchArtistsResultItem {

		i := SearchArtistsResultItem{
			ID:         m.ID,
			Department: m.Department,
			Name:       m.Name,
			Gender:     m.Gender,
		}

		movies := lo.Map(m.Movies, func(m tmdb.Movie, _ int) SearchMoviesResultItem {

			i := SearchMoviesResultItem{
				ID:          m.ID,
				Title:       m.Title,
				ReleaseDate: m.ReleaseDate,
				Synopsis:    m.Overview,
				Genres:      m.Genres,
				Language:    m.Language,
			}

			return i
		})

		i.Movies = movies

		return i
	})

	response := SearchArtistsResponse{
		Results: resultItems,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return []byte(`"error": "failed to marshal response"`), fmt.Errorf("failed to marshal response: %w", err)
	}

	return responseBytes, nil
}
