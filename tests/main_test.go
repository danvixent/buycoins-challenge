package tests

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/danvixent/buycoins_challenge/graphql"
	"github.com/danvixent/buycoins_challenge/handlers/margin"
)

const (
	port    = "8081"
	baseURL = "http://localhost:" + port
)

func TestMain(m *testing.M) {
	marginHandler := margin.NewHandler()
	graphqlHandler := graphql.NewHandler(marginHandler)

	mux := http.NewServeMux()
	graphqlHandler.SetupRoutes(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("serving graphql endpoint at http://localhost:%s/graphql", port)
	log.Printf("serving graphiql endpoint at http://localhost:%s/graphiql", port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("unable to listen: %s", err)
		}
	}()

	// allow the goroutine above start the server
	time.Sleep(time.Second)

	// run the tests
	code := m.Run()

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("unable to shutdown server gracefully: %v", err)
	}

	os.Exit(code)
}
