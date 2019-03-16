CREATE TABLE passers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    campus_id INTEGER REFERENCES campuses(id),
    course_id INTEGER REFERENCES courses(id)
)