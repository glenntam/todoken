
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    created TIMESTAMP NOT NULL,
    email TEXT NOT NULL,
    hashed_password TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_users_lower_email ON users(LOWER(email));

