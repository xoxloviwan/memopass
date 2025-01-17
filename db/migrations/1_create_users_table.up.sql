CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(30) NOT NULL UNIQUE,
    password VARCHAR(72) NOT NULL);
CREATE INDEX users_login_idx ON users (username);