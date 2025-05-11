package movie

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/sparkymat/moviecp/internal/provider/tmdb"
)

var ErrNotFound = errors.New("not found")

type CastItem struct {
	Actor         string `json:"actor"`
	CharacterName string `json:"characterName"`
}

type FetchMovieResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	ReleaseDate string     `json:"releaseData"`
	Synopsis    string     `json:"synopsis"`
	Genres      []string   `json:"genres"`
	Language    string     `json:"language"`
	Cast        []CastItem `json:"cast"`
}

func (s *Service) FetchMovie(ctx context.Context, title string) ([]byte, error) {
	movieResults, _, err := s.tmdb.SearchMovies(ctx, title, 1)
	if err != nil {
		return []byte(`{"error": "failed to query"}`), fmt.Errorf("failed to query movies: %w", err)
	}

	if len(movieResults) == 0 {
		return []byte(`{"error": "failed to query"}`), fmt.Errorf("failed to query movies: %w", ErrNotFound)
	}

	id := movieResults[0].ID

	movie, credits, err := s.tmdb.FetchMovie(ctx, id)
	if err != nil {
		return []byte(`{"error": "failed to query"}`), fmt.Errorf("failed to query movies: %w", err)
	}

	response := FetchMovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		ReleaseDate: movie.ReleaseDate,
		Synopsis:    movie.Overview,
		Genres:      movie.Genres,
		Language:    movie.Language,
	}

	castItems := lo.Map(credits, func(i tmdb.Credit, _ int) CastItem {
		c := CastItem{
			Actor:         i.Name,
			CharacterName: i.Character,
		}

		return c
	})

	response.Cast = castItems

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return []byte(`"error": "failed to marshal response"`), fmt.Errorf("failed to marshal response: %w", err)
	}

	return responseBytes, nil
}
