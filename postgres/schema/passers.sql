CREATE TABLE passers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    course_id INTEGER REFERENCES courses(id) ON DELETE SET NULL,
    campus_id INTEGER REFERENCES campuses(id)
);