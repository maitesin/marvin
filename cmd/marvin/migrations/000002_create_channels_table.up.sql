CREATE TABLE IF NOT EXISTS channels (
    identifier UUID PRIMARY KEY,
    delivery_id TEXT REFERENCES deliveries(identifier)
);