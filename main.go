package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danvixent/buycoins_challenge/graphql"
	"github.com/danvixent/buycoins_challenge/handlers/margin"
)

func main() {
	marginHandler := margin.NewHandler()
	graphqlHandler := graphql.NewHandler(marginHandler)

	mux := http.NewServeMux()
	graphqlHandler.SetupRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("serving graphql endpoint at http://localhost:%s/graphql", port)
	log.Printf("serving graphiql endpoint at http://localhost:%s/graphiql", port)

	// start server in new goroutine so we can listen for CLI signals
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("unable to listen: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so no need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Panic("Server Shutdown: err=%v", err)
	}
	select {
	case <-ctx.Done():
		log.Print("timeout of 1 seconds.")
	}
	log.Print("Server exiting")
}
