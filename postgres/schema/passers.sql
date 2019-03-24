CREATE TABLE passers (
    id SERIAL PRIMARY KEY,
    date_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    course_id INTEGER REFERENCES courses(id),
    campus_id INTEGER REFERENCES campuses(id)
);