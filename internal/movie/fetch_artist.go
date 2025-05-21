package movie

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"github.com/sparkymat/moviecp/internal/provider/tmdb"
)

type FetchArtistResponse struct {
	ID         string                   `json:"id"`
	Department string                   `json:"department"`
	Name       string                   `json:"name"`
	Gender     string                   `json:"gender"`
	CastOf     []SearchMoviesResultItem `json:"castOf"`
	CrewOf     []SearchMoviesResultItem `json:"crewOf"`
}

func (s *Service) FetchArtist(ctx context.Context, name string) ([]byte, error) {
	artists, _, err := s.tmdb.SearchArtists(ctx, name, 1)
	if err != nil {
		return []byte(`{"error": "failed to query"}`), fmt.Errorf("failed to query artists: %w", err)
	}

	if len(artists) == 0 {
		return []byte(`{"error": "failed to query"}`), fmt.Errorf("failed to query artists: %w", ErrNotFound)
	}

	id := artists[0].ID

	response := FetchArtistResponse{
		ID:         artists[0].ID,
		Department: artists[0].Department,
		Name:       artists[0].Name,
		Gender:     artists[0].Gender,
	}

	castMovies, crewMovies, err := s.tmdb.FetchArtistMovies(ctx, id)
	if err != nil {
		return []byte(`{"error": "failed to query"}`), fmt.Errorf("failed to fetch artist movies: %w", err)
	}

	response.CastOf = lo.Map(castMovies, func(m tmdb.Movie, _ int) SearchMoviesResultItem {

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

	response.CrewOf = lo.Map(crewMovies, func(m tmdb.Movie, _ int) SearchMoviesResultItem {

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

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return []byte(`"error": "failed to marshal response"`), fmt.Errorf("failed to marshal response: %w", err)
	}

	return responseBytes, nil
}
