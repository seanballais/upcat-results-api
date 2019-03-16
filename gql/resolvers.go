package gql

import (
    "github.com/seanballais/upcat-results-api/postgres"
    "github.com/graphql-go/graphql"
)

type Resolver struct {
    db *postgres.Db
}

// PasserResolver resolves our passer query through a DB call to GetPassers.
func (r *Resolver) PasserResolver(p graphql.ResolveParams) (interface{}, error) {
    name, ok := p.Args["name"].(string)
    if !ok {
        name = ""
    }

    course, ok := p.Args["course"].(string)
    if !ok {
        course = ""
    }

    campus, ok := p.Args["campus"].(string)
    if !ok {
        campus = ""
    }

    passers := r.db.GetPassers(name, course, campus)
    return passers, nil
}

// CourseResolver resolves our course query through a DB call to GetCourses.
func (r *Resolver) CourseResolver(p graphql.ResolveParams) (interface{}, error) {
    courses := r.db.GetCourses()
    return courses, nil
}