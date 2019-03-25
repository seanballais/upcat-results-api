CREATE TABLE searchRequests (
    id SERIAL PRIMARY KEY,
    date_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    course_id INTEGER REFERENCES courses(id),
    campus_id INTEGER REFERENCES campuses(id),
    location_id TEXT REFERENCES locations(id),
    location_computed_via_gps BOOLEAN NOT NULL
);
