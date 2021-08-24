package http

import (
	"net/http"

	"github.com/rafaelrubbioli/fileapi/pkg/graphql"
	"github.com/rafaelrubbioli/fileapi/pkg/graphql/explorer"
	"github.com/rafaelrubbioli/fileapi/pkg/middleware"
	"github.com/rafaelrubbioli/fileapi/pkg/service"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
)

func NewServer(service service.Service) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(chimiddleware.DefaultLogger)
	graphqlHandler := graphql.NewHandler(service)
	r.With(middleware.CorsMiddleware).
		Route("/graphql", func(r chi.Router) {
			r.Handle("/", http.HandlerFunc(graphqlHandler.Handle))
			r.Get("/explorer", explorer.Handler)
		})

	return r, nil
}
