package main

import (
	"context"
	"net/http"

	"github.com/integrii/flaggy"
	mcpgo "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/redis/go-redis/v9"
	"github.com/sparkymat/moviecp/internal/movie"
	"github.com/sparkymat/moviecp/internal/provider/tmdb"
)

var apiToken = ""

var redisURL = ""

type MovieFetchParams struct {
	Title string `json:"title" jsonschema:"required,description=Title of movie to be fetched"`
}

func main() {
	flaggy.String(&apiToken, "t", "token", "TMDB API Token")
	flaggy.String(&redisURL, "r", "redis", "Redis URL")

	flaggy.Parse()

	if apiToken == "" {
		panic("missing api token")
	}

	if redisURL == "" {
		panic("missing redis URL")
	}

	done := make(chan struct{})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	tmdbProvider := tmdb.New(http.DefaultClient, redisClient, apiToken)

	movieService := movie.New(tmdbProvider)

	server := mcpgo.NewServer(stdio.NewStdioServerTransport())

	err := server.RegisterTool("fetch_movie_by_title", "Fetch movie details by title", func(arguments MovieFetchParams) (*mcpgo.ToolResponse, error) {
		result, searchErr := movieService.FetchMovie(context.Background(), arguments.Title)
		if searchErr != nil {
			panic(searchErr)
		}

		return &mcpgo.ToolResponse{
			Content: []*mcpgo.Content{
				{
					Type: mcpgo.ContentTypeText,
					TextContent: &mcpgo.TextContent{
						Text: string(result),
					},
				},
			},
		}, nil
	})
	if err != nil {
		panic(err)
	}

	err = server.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}
