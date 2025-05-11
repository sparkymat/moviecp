package tmdb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/samber/lo"
	"github.com/sparkymat/moviecp/internal/provider"
)

//nolint:tagliatelle
type fetchMovieResponse struct {
	ID               int64                         `json:"id"`
	Title            string                        `json:"title"`
	Overview         string                        `json:"overview"`
	ReleaseDate      string                        `json:"release_date"`
	Genres           []fetchMovieResponseGenreItem `json:"genres"`
	PosterPath       string                        `json:"poster_path"`
	VoteAverage      float64                       `json:"vote_average"`
	OriginalLanguage string                        `json:"original_language"`
}

type fetchMovieResponseGenreItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type fetchCreditsResponse struct {
	Cast []fetchCreditsResponseCreditItem `json:"cast"`
}

//nolint:tagliatelle
type fetchCreditsResponseCreditItem struct {
	Character   string `json:"character"`
	Name        string `json:"name"`
	ID          int64  `json:"id"`
	Order       int64  `json:"order"`
	ProfilePath string `json:"profile_path"`
}

func (p *Provider) FetchMovie(ctx context.Context, movieID string) (Movie, []Credit, error) {
	url := fmt.Sprintf("%s/movie/%s?api_key=%s", baseURL, movieID, p.apiToken)

	movieItem, err := provider.FetchJSON[fetchMovieResponse](ctx, p.redisClient, p.httpClient, url)
	if err != nil {
		return Movie{}, nil, fmt.Errorf("failed to fetch movies: %w", err)
	}

	url = fmt.Sprintf("%s/movie/%s/credits?api_key=%s", baseURL, movieID, p.apiToken)

	creditItems, err := provider.FetchJSON[fetchCreditsResponse](ctx, p.redisClient, p.httpClient, url)
	if err != nil {
		return Movie{}, nil, fmt.Errorf("failed to fetch credits: %w", err)
	}

	credits := lo.Map(creditItems.Cast, func(creditItem fetchCreditsResponseCreditItem, _ int) Credit {
		return Credit{
			ID:        strconv.FormatInt(creditItem.ID, 10),
			Name:      creditItem.Name,
			Character: creditItem.Character,
			Order:     creditItem.Order,
			ProfilePath: func() *string {
				if creditItem.ProfilePath == "" {
					return nil
				}

				v := fmt.Sprintf("%s%s", imageBaseURL, creditItem.ProfilePath)

				return &v
			}(),
		}
	})

	movie := Movie{
		ID:          strconv.FormatInt(movieItem.ID, 10),
		Title:       movieItem.Title,
		Genres:      lo.Map(movieItem.Genres, func(genre fetchMovieResponseGenreItem, _ int) string { return genre.Name }),
		Overview:    movieItem.Overview,
		PosterPath:  fmt.Sprintf("%s%s", imageBaseURL, movieItem.PosterPath),
		ReleaseDate: movieItem.ReleaseDate,
		Language:    movieItem.OriginalLanguage,
		VoteAverage: movieItem.VoteAverage,
	}

	return movie, credits, nil
}
