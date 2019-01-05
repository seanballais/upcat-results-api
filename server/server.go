package server

import (
    "encoding/json"
    "net/http"

    "github.com/seanballais/upcat-results-api/gql"
    "github.com/go-chi/render"
    "github.com/graphql-go/graphql"
)

type Server struct {
    GqlSchema *graphql.Schema
}

type reqBody struct {
    Query string `json:"query"`
}

func (s *Server) GraphQL() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check to make sure the query was provided in the request body.
        if r.Body == nil {
            http.Error(w, "Must provide GraphQL query in request body.", 400)
            return
        }

        var rBody reqBody

        // Decode the request body into rBody.
        err := json.NewDecoder(r.Body).Decode(&rBody)
        if err != nil {
            http.Error(w, "Error parsing JSON request body.", 400)
        }

        result := gql.ExecuteQuery(rBody.Query, *s.GqlSchema)
        render.JSON(w, r, result)
    }
}