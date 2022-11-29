CREATE TABLE IF NOT EXISTS verifications(
    id SERIAL PRIMARY KEY, uuid UUID,
    kind VARCHAR(20),
    description VARCHAR,
    decline_reason VARCHAR,
    status VARCHAR(20),
    created_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
