package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	log.Fatal(http.ListenAndServe(":9000", router))
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

	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.StripSlashes,
		middleware.Recoverer,
	)

	router.Post("/graphql", gqlServer.GraphQL())

	return router, db
}
