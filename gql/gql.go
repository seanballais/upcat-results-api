package gql

import (
    "fmt"

    "github.com/graphql-go/graphql"
)

// ExecuteQuery runs our GraphQL queries.
func ExecuteQuery(params graphql.Params) *graphql.Result {
    result := graphql.Do(params)

    if len(result.Errors) > 0 {
        fmt.Printf("Unexpected errors inside ExecuteQuery: %v", result.Errors)
    }

    return result
}