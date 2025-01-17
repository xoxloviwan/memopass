CREATE TABLE cards (
  user_id INTEGER NOT NULL,
  ccn TEXT NOT NULL,
  exp TEXT NOT NULL,
  cvv TEXT NOT NULL,
  date DATETIME NOT NULL,
  meta TEXT
);
CREATE INDEX cards_user_id_idx ON cards (user_id);