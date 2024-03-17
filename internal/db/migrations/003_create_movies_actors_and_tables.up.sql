-- +goose Up
CREATE TABLE movies (
                        id SERIAL PRIMARY KEY,
                        title VARCHAR(150) NOT NULL UNIQUE CHECK (char_length(title) >= 1),
                        description VARCHAR(1000),
                        release_date DATE,
                        rating INT CHECK (rating >= 0 AND rating <= 10)
);

CREATE TABLE actors (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(100) NOT NULL UNIQUE ,
                        gender VARCHAR(10) CHECK (gender IN ('male', 'female')),
                        birth_date DATE
);

CREATE TABLE movies_actors (
                               movie_id INT,
                               actor_id INT,
                               PRIMARY KEY (movie_id, actor_id),
                               FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
                               FOREIGN KEY (actor_id) REFERENCES actors(id) ON DELETE CASCADE
);