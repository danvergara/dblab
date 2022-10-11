CREATE TABLE myschema.users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE
);
