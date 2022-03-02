CREATE TABLE IF NOT EXISTS deliveries (
  identifier TEXT PRIMARY KEY,
  events JSONB NOT NULL,
  delivered bool DEFAULT false NOT NULL
);