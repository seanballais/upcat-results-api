package server

import (
    "fmt"
    "encoding/json"
    "net"
    "net/http"

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

        userIPAddress, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            fmt.Printf("%q is not in the format: <IP address>:<port>", r.RemoteAddr)
        }

        // TODO: Preprocess userGPS location, or send an error message
        //       when the GPS location is sent in the wrong format. The
        //       format should be: (latitude, longitude).

        rootValue := map[string]interface{} {
            // Adding in `response` and `request` so that `userLocation` does
            // not feel lonely. Additionally, this is to make this look more
            // akin to a standard HTTP header.
            "response":         w,
            "request":          r,
            "userGPSLocation":  r.Header.Get("SNB-User-GPS-Location"),
            "userIPAddress":    userIPAddress,
        }

        params := graphql.Params{
                Schema:         *s.GqlSchema,
                RequestString:  rBody.Query,
                RootObject:     rootValue,
        }
        result := graphql.Do(params)

        if len(result.Errors) > 0 {
            fmt.Printf("Unexpected errors inside ExecuteQuery: %v", result.Errors)
        }

        render.JSON(w, r, result)
    }
}