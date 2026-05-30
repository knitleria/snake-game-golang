package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"snake_golang/internal/leaderboard"

	"github.com/aws/aws-lambda-go/events"
)

func main() {
	addr := os.Getenv("LEADERBOARD_LOCAL_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	handler, err := leaderboard.NewHandler(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(io.LimitReader(r.Body, 32*1024))
		if err != nil {
			http.Error(w, "bad request body", http.StatusBadRequest)
			return
		}

		query := map[string]string{}
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				query[key] = values[0]
			}
		}

		headers := map[string]string{}
		for key, values := range r.Header {
			if len(values) > 0 {
				headers[key] = values[0]
			}
		}

		resp, err := handler.Handle(r.Context(), events.APIGatewayV2HTTPRequest{
			RawPath:               r.URL.Path,
			RawQueryString:        r.URL.RawQuery,
			Headers:               headers,
			QueryStringParameters: query,
			Body:                  string(body),
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method:    r.Method,
					Path:      r.URL.Path,
					SourceIP:  r.RemoteAddr,
					UserAgent: r.UserAgent(),
				},
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for key, value := range resp.Headers {
			w.Header().Set(key, value)
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write([]byte(resp.Body))
	})

	log.Printf("leaderboard local api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
