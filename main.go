package main

import (
	"encoding/json"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

type MovieSearchparams struct {
	Query string `json:"query" jsonschema:"required,description=Text to search for in movie titles"`
}

func main() {
	done := make(chan struct{})
	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	err := server.RegisterTool("search_movies_by_title", "Search movies by title", func(arguments MovieSearchparams) (*mcp_golang.ToolResponse, error) {
		resp := map[string]any{
			"results": []map[string]any{
				{
					"title":    "Ustad Hotel",
					"year":     2012,
					"synopsis": "Faisi wants to go to UK to become a professional chef but circumstances force him to assist his grandfather in a small restaurant in Kozhikode city, changing his outlook on life forever.",
					"genres":   []string{"Comedy", "Drama"},
				},
				{
					"title":    "Hotel California",
					"year":     2013,
					"synopsis": "Amidst a violent confrontation and illicit affairs, a shadowy criminal operation involving counterfeit DVDs and a notorious don, “Airport Jimmy,” begins to unfold.  Parallel storylines of desire, mystery, and media scrutiny converge around this central criminal enterprise.",
					"genres":   []string{"Comedy"},
				},
			},
			"count": 2,
		}

		respBytes, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}

		return &mcp_golang.ToolResponse{
			Content: []*mcp_golang.Content{
				{
					Type: mcp_golang.ContentTypeText,
					TextContent: &mcp_golang.TextContent{
						Text: string(respBytes),
					},
				},
			},
		}, nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterTool("search_movies_by_genre", "Search movies by genre", func(arguments MovieSearchparams) (*mcp_golang.ToolResponse, error) {
		resp := map[string]any{
			"results": []map[string]any{
				{
					"title":    "Pulival Kalyanam",
					"year":     2012,
					"synopsis": "Faisi wants to go to UK to become a professional chef but circumstances force him to assist his grandfather in a small restaurant in Kozhikode city, changing his outlook on life forever.",
				},
				{
					"title":    "Salt n Pepper",
					"year":     2013,
					"synopsis": "Amidst a violent confrontation and illicit affairs, a shadowy criminal operation involving counterfeit DVDs and a notorious don, “Airport Jimmy,” begins to unfold.  Parallel storylines of desire, mystery, and media scrutiny converge around this central criminal enterprise.",
				},
			},
			"count": 2,
		}

		respBytes, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}

		return &mcp_golang.ToolResponse{
			Content: []*mcp_golang.Content{
				{
					Type: mcp_golang.ContentTypeText,
					TextContent: &mcp_golang.TextContent{
						Text: string(respBytes),
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
