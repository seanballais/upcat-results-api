package postgres

import (
    "fmt"
    "database/sql"

    _ "github.com/lib/pq"

    "github.com/seanballais/upcat-results-api/config"
)

type Db struct {
    *sql.DB
}

// New makes a new database using the connection string, and returns it.
// Otherwise, it returns an error.
func New(connString string) (*Db, error) {
    db, err := sql.Open("postgres", connString)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return &Db{db}, nil
}

// ConnString returns a connection string based on the parameters given in the config file (config.yaml).
func CreateConnectionString() string {
    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        config.DbHost, config.DbPort, config.DbUsername, config.DbPassword, config.DbName,
    )
}

type Passer struct {
    ID      int
    Name    string
    Course  *string
    Campus  string
}

// GetPassersByName is called within our passer query for GraphQL.
func (d *Db) GetPassersByName(name string) []Passer {
    stmt, err := d.Prepare("SELECT id, name, course, campus FROM passers WHERE name LIKE '%' || $1 || '%'")
    if err != nil {
        fmt.Println("GetPassersByName Preparation Error: ", err)
    }

    rows, err := stmt.Query(name)
    defer rows.Close()
    if err != nil {
        fmt.Println("GetPassersByName Query Error: ", err)
    }

    var r Passer
    passers := []Passer{}

    for rows.Next() {
        err = rows.Scan(
            &r.ID,
            &r.Name,
            &r.Course,
            &r.Campus,
        )

        if err != nil {
            fmt.Println("Error scanning rows: ", err)
        }

        passers = append(passers, r)
    }

    return passers
}