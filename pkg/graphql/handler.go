//go:generate go run -mod=mod github.com/99designs/gqlgen --verbose

package graphql

import (
	"net/http"

	"github.com/rafaelrubbioli/fileapi/pkg/graphql/gqlgen"
	"github.com/rafaelrubbioli/fileapi/pkg/graphql/introspection"
	"github.com/rafaelrubbioli/fileapi/pkg/graphql/resolver"
	"github.com/rafaelrubbioli/fileapi/pkg/service"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
)

type Handler struct {
	Schema graphql.ExecutableSchema
	handle http.HandlerFunc
}

func NewHandler(service service.Service) Handler {
	schema := gqlgen.NewExecutableSchema(gqlgen.Config{
		Resolvers:  resolver.New(service),
		Directives: gqlgen.DirectiveRoot{},
	})

	server := handler.New(schema)
	server.AddTransport(transport.Options{})
	server.AddTransport(transport.GET{})
	server.AddTransport(transport.POST{})
	server.AddTransport(transport.MultipartForm{})
	server.SetQueryCache(lru.New(1000))
	server.Use(introspection.Introspection{})

	return Handler{
		Schema: schema,
		handle: server.ServeHTTP,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	h.handle(w, r)
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
