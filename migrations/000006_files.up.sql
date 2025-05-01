
CREATE TABLE files (
  id          uuid PRIMARY KEY,
  owner_id    uuid      NOT NULL,
  owner_table text      NOT NULL, 
  url         text      NOT NULL,
  created_at  timestamptz DEFAULT now()
);
CREATE INDEX ON files(owner_table, owner_id);
