package tmdb

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/samber/lo"
	"github.com/sparkymat/moviecp/internal/provider"
)

type searchPeopleKnownForItem struct {
	ID               int64     `json:"id"`
	GenreIDs         []int64   `json:"genre_ids"`
	MediaType        MediaType `json:"media_type"`
	FirstAirDate     string    `json:"first_air_date"`
	ReleaseDate      string    `json:"release_date"`
	Name             string    `json:"name"`
	Title            string    `json:"title"`
	OriginCountry    []string  `json:"origin_country"`
	OriginalLanguage string    `json:"original_language"`
	OriginalName     string    `json:"original_name"`
	OriginalTitle    string    `json:"original_title"`
	Overview         string    `json:"overview"`
	Popularity       float64   `json:"popularity"`
	PosterPath       string    `json:"poster_path"`
	VoteAverage      float64   `json:"vote_average"`
	VoteCount        int64     `json:"vote_count"`
}

type searchPeopleResponseItem struct {
	ID                 int64                      `json:"id"`
	Gender             int64                      `json:"gender"`
	KnownFor           []searchPeopleKnownForItem `json:"known_for"`
	KnownForDepartment string                     `json:"known_for_department"`
	Name               string                     `json:"name"`
	OriginalName       string                     `json:"original_name"`
	Popularity         float64                    `json:"popularity"`
	ProfilePath        string                     `json:"profile_path"`
}

type searchPeopleResponse struct {
	Results    []searchPeopleResponseItem `json:"results"`
	TotalPages int64                      `json:"total_pages"`
}

func (p *Provider) SearchArtists(ctx context.Context, query string, page int64) ([]Artist, int64, error) {
	searchURL := fmt.Sprintf("%s/search/person?api_key=%s&query=%s&page=%d", baseURL, p.apiToken, url.QueryEscape(query), page)

	response, err := provider.FetchJSON[searchPeopleResponse](ctx, p.redisClient, p.httpClient, searchURL)
	if err != nil {
		return []Artist{}, 0, fmt.Errorf("failed to fetch movies: %w", err)
	}

	genresMap, err := p.ListGenres(ctx, "en") // Get genres names in en
	if err != nil {
		return []Artist{}, 0, fmt.Errorf("failed to fetch genres: %w", err)
	}

	artists := lo.Map(response.Results, func(item searchPeopleResponseItem, _ int) Artist {

		a := Artist{
			ID:          strconv.FormatInt(item.ID, 10),
			Department:  item.KnownForDepartment,
			Name:        item.Name,
			Popularity:  item.Popularity,
			ProfilePath: item.ProfilePath,
		}

		switch item.Gender {
		case 1:
			a.Gender = "Male"
		case 2:
			a.Gender = "Female"
		case 0:
			a.Gender = "Unknown"
		default:
			a.Gender = "Unknown"
		}

		for _, knownFor := range item.KnownFor {
			genres := lo.Map(knownFor.GenreIDs, func(genreID int64, _ int) string { return genresMap[genreID] })
			switch knownFor.MediaType {
			case MediaTypeMovie:
				movie := Movie{
					ID:          strconv.FormatInt(knownFor.ID, 10),
					Title:       knownFor.Title,
					Genres:      genres,
					Overview:    knownFor.Overview,
					ReleaseDate: knownFor.ReleaseDate,
					Language:    knownFor.OriginalLanguage,
					PosterPath:  knownFor.PosterPath,
					VoteAverage: knownFor.VoteAverage,
				}

				a.Movies = append(a.Movies, movie)
			case MediaTypeTvShow:
				tvShow := TvShow{
					ID:          strconv.FormatInt(knownFor.ID, 10),
					Title:       knownFor.Name,
					Genres:      genres,
					Overview:    knownFor.Overview,
					Language:    knownFor.OriginalLanguage,
					PosterPath:  knownFor.PosterPath,
					VoteAverage: knownFor.VoteAverage,
				}

				a.TvShows = append(a.TvShows, tvShow)
			default:
			}
		}

		return a
	})

	return artists, response.TotalPages, nil
}
