package movie

import (
	"context"

	"github.com/sparkymat/moviecp/internal/provider/tmdb"
)

type TMDBProvider interface {
	SearchMovies(ctx context.Context, query string, page int64) ([]tmdb.Movie, int64, error)
}
