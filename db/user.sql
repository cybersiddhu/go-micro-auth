CREATE TABLE users (
     id INTEGER PRIMARY KEY autoincrement not null,
     email TEXT UNIQUE not null,
     password TEXT UNIQUE not null,
     firstname TEXT,
     lastname TEXT,
     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
     updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
