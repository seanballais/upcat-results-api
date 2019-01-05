package gql

import (
    "github.com/seanballais/upcat-results-api/postgres"
    "github.com/graphql-go/graphql"
)

type Resolver struct {
    db *postgres.Db
}

// PasserResolver resolves our passer query through a DB call to
// GetPassersByName.
func (r *Resolver) PasserResolver(p graphql.ResolveParams) (interface{}, error) {
    name, ok := p.Args["name"].(string)
    if ok {
        passers := r.db.GetPassersByName(name)
        return passers, nil
    }

    return nil, nil
}