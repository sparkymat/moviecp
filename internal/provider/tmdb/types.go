//nolint:tagliatelle
package tmdb

type MediaType string

const (
	MediaTypeMovie  MediaType = "movie"
	MediaTypeTvShow MediaType = "tv"
)

type Movie struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Genres      []string `json:"genres"`
	Overview    string   `json:"overview"`
	PosterPath  string   `json:"poster_path"`
	ReleaseDate string   `json:"release_date"`
	Language    string   `json:"original_language"`
	VoteAverage float64  `json:"vote_average"`
}

type Artist struct {
	ID          string   `json:"id"`
	Department  string   `json:"department"`
	Name        string   `json:"name"`
	Popularity  float64  `json:"popularity"`
	ProfilePath string   `json:"profile_path"`
	Gender      string   `json:"gender"`
	Movies      []Movie  `json:"movie_ids"`
	TvShows     []TvShow `json:"tv_show_ids"`
}

type Crew struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	ProfilePath *string `json:"profile_path"`
	Department  string  `json:"department"`
	Job         string  `json:"job"`
}

type Cast struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Character   string  `json:"character"`
	Order       int64   `json:"order"`
	ProfilePath *string `json:"profile_path"`
}

type TvShow struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Genres      []string `json:"genres"`
	Overview    string   `json:"overview"`
	PosterPath  string   `json:"poster_path"`
	Language    string   `json:"original_language"`
	VoteAverage float64  `json:"vote_average"`
}

type Season struct {
	Number     int64     `json:"number"`
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	PosterPath string    `json:"poster_path"`
	Episodes   []Episode `json:"episodes"`
	Credits    []Cast    `json:"credits"`
}

type Episode struct {
	Number     int64  `json:"number"`
	Name       string `json:"name"`
	AirDate    string `json:"air_date"`
	ID         string `json:"id"`
	Overview   string `json:"overview"`
	GuestStars []Cast `json:"guest_stars"`
}
