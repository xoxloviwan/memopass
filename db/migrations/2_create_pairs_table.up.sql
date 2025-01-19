CREATE TABLE pairs (
  id  INTEGER PRIMARY KEY AUTOINCREMENT,
  login TEXT NOT NULL,
  password TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  date DATETIME NOT NULL,
  meta TEXT
);
CREATE INDEX pairs_user_id_idx ON pairs (user_id);