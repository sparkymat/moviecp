package movie

func New(tmdbProvider TMDBProvider) *Service {
	return &Service{
		tmdb: tmdbProvider,
	}
}

type Service struct {
	tmdb TMDBProvider
}
