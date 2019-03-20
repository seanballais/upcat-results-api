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

type Course struct {
    ID      int
    Name    string
}

type Campus struct {
    ID      int
    Name    string
}

type PassersMetadata struct {
    Num_items int
}

// GetPassers is called within our passer query for GraphQL.
func (d *Db) GetPassers(name string, course_id int, campus_id int, page_number int) []Passer {
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

    if course_id != 0 {
        clause := fmt.Sprintf("passers.course_id=$%d", parameterCtr)
        whereClauses = append(whereClauses, clause)
        parameters = append(parameters, course_id)

        parameterCtr++
    }

    if campus_id != 0 {
        clause := fmt.Sprintf("passers.campus_id=$%d", parameterCtr)
        whereClauses = append(whereClauses, clause)
        parameters = append(parameters, campus_id)

        parameterCtr++
    }

    whereClause := strings.Join(whereClauses, " AND ")
    if whereClause != "" {
        query = fmt.Sprintf("%s AND %s", query, whereClause)
    }

    query += " ORDER BY passers.name ASC"

    passer_page_size := 10  // Arbitrarily chosen 10, cause why not?
    passer_page_start_index := page_number * passer_page_size

    query += fmt.Sprintf(
        " LIMIT %d OFFSET %d",
        passer_page_size,
        passer_page_start_index,
    )

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

    var p Passer
    passers := []Passer{}
    for rows.Next() {
        err = rows.Scan(
            &p.ID,
            &p.Name,
            &p.Course,
            &p.Campus,
        )

        if err != nil {
            fmt.Println("Error scanning Passer rows: ", err)
        }

        passers = append(passers, p)
    }

    return passers
}

func (d *Db) GetCourses() []Course {
    query := "SELECT id, name FROM courses"

    stmt, err := d.Prepare(query)
    if err != nil {
        fmt.Println("GetCourses Preparation Error: ", err)
    }

    rows, err := stmt.Query()
    defer rows.Close()
    if err != nil {
        fmt.Println("GetCourses Query Error: ", err)
    }

    var c Course
    courses := []Course{}
    for rows.Next() {
        err = rows.Scan(
            &c.ID,
            &c.Name,
        )

        if err != nil {
            fmt.Println("Error scanning Course rows: ", err)
        }

        courses = append(courses, c)
    }

    return courses
}

func (d *Db) GetCampuses() []Campus {
    query := "SELECT id, name FROM campuses"

    stmt, err := d.Prepare(query)
    if err != nil {
        fmt.Println("GetCampuses Preparation Error: ", err)
    }

    rows, err := stmt.Query()
    defer rows.Close()
    if err != nil {
        fmt.Println("GetCampuses Query Error: ", err)
    }

    var c Campus
    campuses := []Campus{}
    for rows.Next() {
        err = rows.Scan(
            &c.ID,
            &c.Name,
        )

        if err != nil {
            fmt.Println("Error scanning Campus rows: ", err)
        }

        campuses = append(campuses, c)
    }

    return campuses
}

func (d *Db) GetPassersMetadata(name string, course_id int, campus_id int) PassersMetadata {
    query := "SELECT COUNT(*)"
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

    if course_id != 0 {
        clause := fmt.Sprintf("passers.course_id=$%d", parameterCtr)
        whereClauses = append(whereClauses, clause)
        parameters = append(parameters, course_id)

        parameterCtr++
    }

    if campus_id != 0 {
        clause := fmt.Sprintf("passers.campus_id=$%d", parameterCtr)
        whereClauses = append(whereClauses, clause)
        parameters = append(parameters, campus_id)

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

    var passersMetadata PassersMetadata
    rows.Next()
    err = rows.Scan(
        &passersMetadata.Num_items,
    )

    if err != nil {
        fmt.Println("Error scanning PassersMetadata row: ", err)
    }

    return passersMetadata
}
