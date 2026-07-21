CREATE TABLE IF NOT EXISTS products (
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    price      DOUBLE PRECISION NOT NULL,
    stock      INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);
