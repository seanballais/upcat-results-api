CREATE TABLE searchRequests (
    id SERIAL PRIMARY KEY,
    date_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    course_id INTEGER REFERENCES courses(id) ON DELETE SET NULL,
    campus_id INTEGER REFERENCES campuses(id),
    location TEXT,
    location_computed_via_gps BOOLEAN NOT NULL,
    src_ip_address VARCHAR(15) NOT NULL
);
