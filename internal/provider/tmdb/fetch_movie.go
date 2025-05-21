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
	ID                  int64                                 `json:"id"`
	Title               string                                `json:"title"`
	Overview            string                                `json:"overview"`
	ReleaseDate         string                                `json:"release_date"`
	Genres              []fetchMovieResponseGenreItem         `json:"genres"`
	PosterPath          string                                `json:"poster_path"`
	VoteAverage         float64                               `json:"vote_average"`
	OriginalLanguage    string                                `json:"original_language"`
	ProductionCompanies []fetchMovieResponseProductionCompany `json:"production_companies"`
}

type fetchMovieResponseProductionCompany struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
	LogoPath      string `json:"logo_path"`
}

type fetchMovieResponseGenreItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type fetchCreditsResponse struct {
	Cast []fetchCreditsResponseCastItem `json:"cast"`
	Crew []fetchCreditsResponseCrewItem `json:"crew"`
}

//nolint:tagliatelle
type fetchCreditsResponseCastItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ProfilePath string `json:"profile_path"`
	Gender      int64  `json:"gender"`
	Character   string `json:"character"`
	Order       int64  `json:"order"`
}

type fetchCreditsResponseCrewItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ProfilePath string `json:"profile_path"`
	Gender      int64  `json:"gender"`
	Department  string `json:"department"`
	Job         string `json:"job"`
}

func (p *Provider) FetchMovie(ctx context.Context, movieID string) (Movie, []Cast, []Crew, error) {
	url := fmt.Sprintf("%s/movie/%s?api_key=%s", baseURL, movieID, p.apiToken)

	movieItem, err := provider.FetchJSON[fetchMovieResponse](ctx, p.redisClient, p.httpClient, url)
	if err != nil {
		return Movie{}, nil, nil, fmt.Errorf("failed to fetch movies: %w", err)
	}

	url = fmt.Sprintf("%s/movie/%s/credits?api_key=%s", baseURL, movieID, p.apiToken)

	creditItems, err := provider.FetchJSON[fetchCreditsResponse](ctx, p.redisClient, p.httpClient, url)
	if err != nil {
		return Movie{}, nil, nil, fmt.Errorf("failed to fetch credits: %w", err)
	}

	cast := lo.Map(creditItems.Cast, func(creditItem fetchCreditsResponseCastItem, _ int) Cast {
		c := Cast{
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

		switch creditItem.Gender {
		case 0:
			c.Gender = "Unknown"
		case 1:
			c.Gender = "Female"
		case 2:
			c.Gender = "Male"
		}

		return c
	})

	crew := lo.Map(creditItems.Crew, func(creditItem fetchCreditsResponseCrewItem, _ int) Crew {
		c := Crew{
			ID:   strconv.FormatInt(creditItem.ID, 10),
			Name: creditItem.Name,
			ProfilePath: func() *string {
				if creditItem.ProfilePath == "" {
					return nil
				}

				v := fmt.Sprintf("%s%s", imageBaseURL, creditItem.ProfilePath)

				return &v
			}(),
			Department: creditItem.Department,
			Job:        creditItem.Job,
		}

		switch creditItem.Gender {
		case 0:
			c.Gender = "Unknown"
		case 1:
			c.Gender = "Female"
		case 2:
			c.Gender = "Male"
		}

		return c
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

	return movie, cast, crew, nil
}
