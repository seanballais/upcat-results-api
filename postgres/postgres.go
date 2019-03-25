package postgres

import (
    "os"
    "fmt"
    "strings"
    "strconv"
    "database/sql"

    _ "github.com/lib/pq"
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
    dbName := os.Getenv("UPCAT_RESULTS_API_DB_NAME")
    dbHost := os.Getenv("UPCAT_RESULTS_API_DB_HOST")
    dbPort, _ := strconv.ParseUint(os.Getenv("UPCAT_RESULTS_API_DB_PORT"), 10, 16)
    dbUsername := os.Getenv("UPCAT_RESULTS_API_DB_USERNAME")
    dbPassword := os.Getenv("UPCAT_RESULTS_API_DB_PASSWORD")

    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s",
        dbHost, dbPort, dbUsername, dbPassword, dbName,
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
    // We should capitalize `name` first because the passers' names
    // are capitalized in the database.
    name = strings.ToUpper(name)

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
    err = rows.Scan(&passersMetadata.Num_items)

    if err != nil {
        fmt.Println("Error scanning PassersMetadata row: ", err)
    }

    return passersMetadata
}

func (d *Db) GetCurrentMonthMappedIPAddresses() int {
    query := "SELECT COUNT(*) FROM ipAddressLocations "
    query += "WHERE date_created == date_trunc('month', CURRENT_DATE)"

    stmt, err := d.Prepare(query)
    if err != nil {
        fmt.Println("GetCurrentMonthMappedIPAddresses Preparation Error: ", err)
    }

    rows, err := stmt.Query()
    defer rows.Close()
    if err != nil {
        fmt.Println("GetCurrentMonthMappedIPAddresses Query Error: ", err)
    }

    var numIPAddresses int
    rows.Next()
    err = rows.Scan(&numIPAddresses)
    if err != nil {
        fmt.Println("Error scanning GetCurrentMonthMappedIPAddresses rows: ", err)
    }

    return numIPAddresses
}

func (d *Db) IsIPAddressCached(ipAddress string) bool {
    query := "SELECT COUNT(*) FROM ipAddressLocations "
    query += "WHERE ip_address=$1"

    stmt, err := d.Prepare(query)
    if err != nil {
        fmt.Println("IsIPAddressCached Preparation Error: ", err)
    }

    rows, err := stmt.Query()
    defer rows.Close()
    if err != nil {
        fmt.Println("IsIPAddressCached Query Error: ", err)
    }

    var numIPAddresses int
    rows.Next()
    err = rows.Scan(&numIPAddresses)
    if err != nil {
        fmt.Println("Error scanning IsIPAddressCached rows: ", err)
    }

    return numIPAddresses > 0
}

func (d *Db) GetIPAddressLocationID(ipAddress string) int {
    query := "SELECT location_id FROM ipAddressLocations "
    query += "WHERE ip_address=$1"

    stmt, err := d.Prepare(query)
    if err != nil {
        fmt.Println("GetIPAddressLocationID Preparation Error: ", err)
    }

    rows, err := stmt.Query()
    defer rows.Close()
    if err != nil {
        fmt.Println("GetIPAddressLocationID Query Error: ", err)
    }

    var locationID int
    rows.Next()
    err = rows.Scan(&locationID)
    if err != nil {
        fmt.Println("Error scanning GetIPAddressLocationID rows: ", err)
    }

    return locationID
}

func (d *Db) AddIPAddressLocationMapping(ip_address string, location_id int) {
    query := "INSERT INTO ipAddressLocations(ip_address, location_id) "
    query += "VALUES($1, $2)"

    stmt, err := d.Prepare(query)
    if err != nil {
        fmt.Println("AddIPAddressLocationMapping Preparation Error: ", err)
    }

    rows, err := stmt.Query(ip_address, location_id)
    defer rows.Close()
    if err != nil {
        fmt.Println("AddIPAddressLocationMapping Query Error: ", err)
    }
}

func (d *Db) AddLocation(location string) int {
    query := "INSERT INTO locations(name) VALUES($1) RETURNING id"

    stmt, err := d.Prepare(query)
    if err != nil {
        fmt.Println("AddLocation Preparation Error: ", err)
    }

    rows, err := stmt.Query(location)
    defer rows.Close()
    if err != nil {
        fmt.Println("AddLocation Query Error: ", err)
    }

    var locationID int
    rows.Next()
    err = rows.Scan(&locationID)
    if err != nil {
        fmt.Println("Error scanning AddLocation rows: ", err)
    }

    return locationID
}

func (d *Db) AddSearchQuery(name string, course_id int, campus_id int, page_number int,
                            location_id int, isLocationComputedViaGPS bool) {
    query := "INSERT INTO searchRequests(name, course_id, campus_id, page_number,"
    query += "location_id, location_computed_via_gps) "
    query += "VALUES($1, $2, $3, $4, $5, $6)"

    var parameters = []interface{}{}

    if name != "" {
        parameters = append(parameters, name)
    } else {
        parameters = append(parameters, nil)
    }

    if course_id != 0 {
        parameters = append(parameters, course_id)
    } else {
        parameters = append(parameters, nil)
    }

    if campus_id != 0 {
        parameters = append(parameters, campus_id)
    } else {
        parameters = append(parameters, nil)
    }

    parameters = append(parameters, page_number)
    parameters = append(parameters, location_id)
    parameters = append(parameters, isLocationComputedViaGPS)

    stmt, err := d.Prepare(query)
    if err != nil {
        fmt.Println("AddSearchQuery Preparation Error: ", err)
    }

    rows, err := stmt.Query(parameters...)
    defer rows.Close()
    if err != nil {
        fmt.Println("AddSearchQuery Query Error: ", err)
    }
}
