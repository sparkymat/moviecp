package movie

import (
	"context"

	"github.com/sparkymat/moviecp/internal/provider/tmdb"
)

type TMDBProvider interface {
	FetchMovie(ctx context.Context, movieID string) (tmdb.Movie, []tmdb.Credit, error)
	SearchArtists(ctx context.Context, query string, page int64) ([]tmdb.Artist, int64, error)
	SearchMovies(ctx context.Context, query string, page int64) ([]tmdb.Movie, int64, error)
}
