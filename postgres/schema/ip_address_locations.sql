CREATE TABLE ipAddressLocations (
    id SERIAL PRIMARY KEY,
    date_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address VARCHAR(15) UNIQUE NOT NULL,
    location_id INTEGER REFERENCES locations(id) NOT NULL
);