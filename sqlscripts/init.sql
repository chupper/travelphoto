CREATE TABLE gallery(
    ID SERIAL PRIMARY KEY,
    NAME TEXT NOT NULL,
    DESCRIPTION TEXT NOT NULL,
    SHOW BOOLEAN NOT NULL,
    RANK INT
);