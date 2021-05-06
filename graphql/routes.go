package graphql

import (
	"log"
	"net/http"

	"github.com/friendsofgo/graphiql"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/danvixent/buycoins_challenge/handlers/margin"
)

type Handler struct {
	marginHandler *margin.Handler
}

const (
	graphqlEndpoint  = "/graphql"
	graphiqlEndpoint = "/graphiql"
)

func NewHandler(marginHandler *margin.Handler) *Handler {
	return &Handler{marginHandler: marginHandler}
}

func (h *Handler) graphqlHandler() http.HandlerFunc {
	c := Config{
		Resolvers: &Resolver{marginHandler: h.marginHandler},
	}

	s := handler.NewDefaultServer(NewExecutableSchema(c))

	return s.ServeHTTP
}

func (h *Handler) graphiqlHandler() http.HandlerFunc {
	graphiqlHandler, err := graphiql.NewGraphiqlHandler(graphqlEndpoint)
	if err != nil {
		log.Panic(err)
	}
	return graphiqlHandler.ServeHTTP
}

func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	graphqlHandlerFunc := h.graphqlHandler()
	mux.HandleFunc(graphqlEndpoint, handleMethod(http.MethodPost, graphqlHandlerFunc))

	graphiqlHandlerFunc := h.graphiqlHandler()
	mux.HandleFunc(graphiqlEndpoint, handleMethod(http.MethodGet, graphiqlHandlerFunc))
}

func handleMethod(method string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusBadGateway)
			return
		}
		handlerFunc(w, r)
	}
}
