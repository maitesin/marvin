CREATE TABLE IF EXISTS deliveries (
  identifier TEXT PRIMARY KEY,
  events JSONB NOT NULL,
  delivered bool DEFAULT false NOT NULL
);