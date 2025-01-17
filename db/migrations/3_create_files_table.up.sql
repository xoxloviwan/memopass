CREATE TABLE files (
  file BLOB NOT NULL,
  name TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  date DATETIME NOT NULL,
  binary BOOLEAN NOT NULL,
  meta TEXT
);
CREATE INDEX files_user_id_idx ON files (user_id);