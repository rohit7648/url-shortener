package server

import (
	"context"
	"encoding/json"
	"net/http"
	v1 "url-shortener/api/urlshortener/v1"
	"url-shortener/internal/conf"
	"url-shortener/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, urlShortener *service.UrlShortenerService, logger log.Logger) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := khttp.NewServer(opts...)

	// Register HTTP handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			LongURL string `json:"long_url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		reply, err := urlShortener.Shorten(context.Background(), &v1.ShortenRequest{
			LongUrl: req.LongURL,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"short_url":  reply.ShortUrl,
			"expires_at": reply.ExpiresAt,
		})
	})

	mux.HandleFunc("/resolve/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		shortCode := r.URL.Path[len("/resolve/"):]
		reply, err := urlShortener.Resolve(context.Background(), &v1.ResolveRequest{ShortCode: shortCode})
		if err != nil {
			if err.Error() == "URL not found or has expired" {
				http.Error(w, "URL not found or has expired", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"long_url":   reply.LongUrl,
			"expires_at": reply.ExpiresAt,
		})
	})

	srv.HandlePrefix("/", mux)
	return srv
}
