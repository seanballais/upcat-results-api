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

    course, ok := p.Args["course"].(int)
    if !ok {
        course = 0
    }

    campus, ok := p.Args["campus"].(int)
    if !ok {
        campus = 0
    }

    page_number, ok := p.Args["page_number"].(int)
    if !ok {
        page_number = 0
    }

    passers := r.db.GetPassers(name, course, campus, page_number)

    // Add the search request to the database.
    rootValue := p.Info.RootValue.(map[string]interface{})
    userLocationID := GetUserLocationID(rootValue["userGPSLocation"].(string),
                                        rootValue["userIPAddress"].(string))

    isLocationComputedViaGPS := true
    if rootValue["userGPSLocation"] == "" {
        isLocationComputedViaGPS = false
    }

    r.db.AddSearchQuery(name, course, campus, page_number, userLocationID, isLocationComputedViaGPS)

    return passers, nil
}

// CourseResolver resolves our course query through a DB call to GetCourses.
func (r *Resolver) CourseResolver(p graphql.ResolveParams) (interface{}, error) {
    courses := r.db.GetCourses()
    return courses, nil
}

// CampusResolver resolves our campus query through a DB call to GetCampuses.
func (r *Resolver) CampusResolver(p graphql.ResolveParams) (interface{}, error) {
    campuses := r.db.GetCampuses()
    return campuses, nil
}

func (r *Resolver) PassersMetadataResolver(p graphql.ResolveParams) (interface{}, error) {
    name, ok := p.Args["name"].(string)
    if !ok {
        name = ""
    }

    course, ok := p.Args["course"].(int)
    if !ok {
        course = 0
    }

    campus, ok := p.Args["campus"].(int)
    if !ok {
        campus = 0
    }

    metadata := r.db.GetPassersMetadata(name, course, campus)
    return metadata, nil
}
