package postgres

import (
    "fmt"
    "strings"
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
func (d *Db) GetPassers(name string, course string, campus string) []Passer {
    query := "SELECT passers.id, passers.name, courses.name, campuses.name "
    query += "FROM passers, courses, campuses "
    query += "WHERE passers.course_id = courses.id "
    query +=   "AND passers.campus_id = campuses.id"

    // Generate the WHERE clause if there are passed parameters.
    var whereClauses = make([]string, 0, 3)
    var parameters = []interface{}{}
    parameterCtr := 1

    if name != "" {
        clause := fmt.Sprintf("passers.name LIKE '%%' || $%d || '%%'",
                              parameterCtr)
        whereClauses = append(whereClauses, clause)
        parameters = append(parameters, name)

        parameterCtr++
    }

    if course != "" {
        clause := fmt.Sprintf("passers.course_id=$%d", parameterCtr)
        whereClauses = append(whereClauses, clause)
        parameters = append(parameters, course)

        parameterCtr++
    }

    if campus != "" {
        clause := fmt.Sprintf("passers.campus_id=$%d", parameterCtr)
        whereClauses = append(whereClauses, clause)
        parameters = append(parameters, campus)

        parameterCtr++
    }

    whereClause := strings.Join(whereClauses, " AND ")
    if whereClause != "" {
        query = fmt.Sprintf("%s AND %s", query, whereClause)
    }

    // Perform SQL query.
    stmt, err := d.Prepare(query)
    if err != nil {
        fmt.Println("GetPassers Preparation Error: ", err)
    }

    rows, err := stmt.Query(parameters...)
    defer rows.Close()
    if err != nil {
        fmt.Println("GetPassers Query Error: ", err)
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