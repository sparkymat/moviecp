package tmdb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/samber/lo"
	"github.com/sparkymat/moviecp/internal/provider"
)

type fetchArtistMoviesResponse struct {
	Cast []searchMoviesResponseItem `json:"cast"`
	Crew []searchMoviesResponseItem `json:"crew"`
}

func (p *Provider) FetchArtistMovies(ctx context.Context, id string) ([]Movie, []Movie, error) {
	searchURL := fmt.Sprintf("%s/person/%s/movie_credits?api_key=%s", baseURL, id, p.apiToken)

	response, err := provider.FetchJSON[fetchArtistMoviesResponse](ctx, p.redisClient, p.httpClient, searchURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch movies: %w", err)
	}

	genresMap, err := p.ListGenres(ctx, "en") // Get genres names in en
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch genres: %w", err)
	}

	castMovies := lo.Map(response.Cast, func(item searchMoviesResponseItem, _ int) Movie {
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

	crewMovies := lo.Map(response.Crew, func(item searchMoviesResponseItem, _ int) Movie {
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

	return castMovies, crewMovies, nil
}
