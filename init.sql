CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS specialty (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS specialist (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    specialty_id INT,
    location GEOGRAPHY,
    address VARCHAR(255),
    url VARCHAR(255),
    telephone VARCHAR(255),
    email VARCHAR(255),
    monday VARCHAR(255),
    tuesday VARCHAR(255),
    wednesday VARCHAR(255),
    thursday VARCHAR(255),
    friday VARCHAR(255),
    saturday VARCHAR(255),
    sunday VARCHAR(255),
    FOREIGN KEY (specialty_id) REFERENCES specialty(id)
);

CREATE TABLE IF NOT EXISTS review (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    specialist_id INT,
    url VARCHAR(255) NOT NULL,
    rating DECIMAL(2, 1) NOT NULL,
    comment VARCHAR(255),
    FOREIGN KEY (specialist_id) REFERENCES specialist(id)
);
