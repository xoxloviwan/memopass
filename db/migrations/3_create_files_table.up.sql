CREATE TABLE files (
  file BLOB NOT NULL,
  name TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  date DATETIME NOT NULL,
  meta TEXT
)