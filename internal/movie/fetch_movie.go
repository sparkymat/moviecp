package movie

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/samber/lo"
	"github.com/sparkymat/moviecp/internal/provider/tmdb"
)

var ErrNotFound = errors.New("not found")

type CastItem struct {
	Name      string `json:"name"`
	Character string `json:"character"`
}

type CrewItem struct {
	Name string `json:"name"`
	Job  string `json:"job"`
}

type FetchMovieResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	ReleaseDate string     `json:"releaseData"`
	Synopsis    string     `json:"synopsis"`
	Genres      []string   `json:"genres"`
	Language    string     `json:"language"`
	Cast        []CastItem `json:"cast"`
	Crew        []CrewItem `json:"crew"`
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

	movie, cast, crew, err := s.tmdb.FetchMovie(ctx, id)
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

	slices.SortFunc(cast, func(a, b tmdb.Cast) int {
		return int(a.Order - b.Order)
	})

	castItems := lo.Map(cast, func(i tmdb.Cast, _ int) CastItem {
		c := CastItem{
			Name:      i.Name,
			Character: i.Character,
		}

		return c
	})

	response.Cast = castItems

	crewItems := lo.Map(crew, func(i tmdb.Crew, _ int) CrewItem {
		c := CrewItem{
			Name: i.Name,
			Job:  i.Job,
		}

		return c
	})

	response.Crew = crewItems

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return []byte(`"error": "failed to marshal response"`), fmt.Errorf("failed to marshal response: %w", err)
	}

	return responseBytes, nil
}
