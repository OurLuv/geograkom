CREATE TABLE IF NOT EXISTS routes (
    route_id SERIAL PRIMARY KEY,
    route_name VARCHAR(255),
    load FLOAT,
    cargo_type VARCHAR(255),
    is_actual BOOLEAN
);