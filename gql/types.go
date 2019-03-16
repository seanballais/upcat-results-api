package gql

import (
    "github.com/graphql-go/graphql"
)

var Passer = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Passer",
        Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.Int,
            },
            "name": &graphql.Field{
                Type: graphql.String,
            },
            "course": &graphql.Field{
                Type: graphql.String,
            },
            "campus": &graphql.Field{
                Type: graphql.String,
            },
        },
    },
)

var Course = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Course",
        Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.Int,
            },
            "name": &graphql.Field{
                Type: graphql.String,
            },
        },
    },
)