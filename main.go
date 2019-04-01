package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
	"github.com/seanballais/upcat-results-api/gql"
	"github.com/seanballais/upcat-results-api/postgres"
	"github.com/seanballais/upcat-results-api/server"
)

func main() {
	fmt.Println("Initialized server.")

	router, db := initializeAPI()
	defer db.Close()

	port := fmt.Sprintf(":%s", os.Getenv("UPCAT_RESULTS_API_PORT"))
	log.Fatal(http.ListenAndServe(port, router))
}

func initializeAPI() (*chi.Mux, *postgres.Db) {
	router := chi.NewRouter()

	db, err := postgres.New(postgres.CreateConnectionString())

	if err != nil {
		log.Fatal(err)
	}

	rootQuery := gql.NewRoot(db)
	sc, err := graphql.NewSchema(graphql.SchemaConfig{Query: rootQuery.Query})

	if err != nil {
		fmt.Println("Error creating schema: ", err)
	}

	gqlServer := server.Server{
		GqlSchema: &sc,
	}

	// Enable CORS.
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "OPTIONS"},
    	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    	ExposedHeaders:   []string{"Link"},
    	AllowCredentials: true,
    	MaxAge:           300,
	})

	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.StripSlashes,
		middleware.Recoverer,
		cors.Handler,
	)

	router.Post("/graphql", gqlServer.GraphQL())

	return router, db
}
