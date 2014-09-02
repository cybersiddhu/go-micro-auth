CREATE TABLE users (
     id serial PRIMARY KEY,
     email varchar(40) NOT NULL,
     password varchar(256) NOT NULL,
     firstname varchar(50),
     lastname varchar(50),
     created_at timestamp DEFAULT CURRENT_TIMESTAMP,
     updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
     UNIQUE(email),
     UNIQUE(password)
);
