package gql

import (
    "github.com/seanballais/upcat-results-api/postgres"
    "github.com/graphql-go/graphql"
)

type Root struct {
    Query *graphql.Object
}

func NewRoot(db *postgres.Db) *Root {
    resolver := Resolver{db: db}

    root := Root{
        Query: graphql.NewObject(
            graphql.ObjectConfig{
                Name: "Query",
                Fields: graphql.Fields{
                    "passers": &graphql.Field{
                        Type: graphql.NewList(Passer),
                        Args: graphql.FieldConfigArgument{
                            "name": &graphql.ArgumentConfig{
                                Type: graphql.String,
                            },
                            "course": &graphql.ArgumentConfig{
                                Type: graphql.Int,
                            },
                            "campus": &graphql.ArgumentConfig{
                                Type: graphql.Int,
                            },
                            "page_number": &graphql.ArgumentConfig{
                                Type: graphql.Int,
                            },
                        },
                        Resolve: resolver.PasserResolver,
                    },
                    "courses": &graphql.Field{
                        Type: graphql.NewList(Course),
                        Resolve: resolver.CourseResolver,
                    },
                    "campuses": &graphql.Field{
                        Type: graphql.NewList(Campus),
                        Resolve: resolver.CampusResolver,                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           
                    },
                    "passers_metadata": &graphql.Field{
                        Type: PassersMetadata,
                        Args: graphql.FieldConfigArgument{
                            "name": &graphql.ArgumentConfig{
                                Type: graphql.String,
                            },
                            "course": &graphql.ArgumentConfig{
                                Type: graphql.Int,
                            },
                            "campus": &graphql.ArgumentConfig{
                                Type: graphql.Int,
                            },
                        },
                        Resolve: resolver.PassersMetadataResolver,
                    },
                },
            },
        ),
    }

    return &root
}