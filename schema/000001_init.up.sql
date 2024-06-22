CREATE SEQUENCE routes_route_id_seq START 1;

CREATE TABLE routes (
    route_id INTEGER DEFAULT nextval('routes_route_id_seq') PRIMARY KEY,
    route_name VARCHAR(255),
    load FLOAT,
    cargo_type VARCHAR(255),
    is_actual BOOLEAN
);