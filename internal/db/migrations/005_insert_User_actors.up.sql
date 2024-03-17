-- +goose Up
INSERT INTO users (login, pass, role_id)
VALUES ('admin', '$2a$10$7hJcPCmyuOgL82Wx5JG/s.T6jWqskCScYTvQEIaFNVY.jqgFMkuSe', 2);

INSERT INTO actors (name, gender, birth_date)
VALUES ('Cillian Murphy', 'male', '1976-05-25'),
       ('Matt Damon', 'male', '1970-10-08'),
       ('Emily Blunt', 'female', '1983-02-23'),
       ('Florence Pugh', 'female', '1996-01-03'),
       ('Zendaya', 'female', '1996-09-01'),
       ('Timothée Chalamet', 'male', '1995-12-27'),
       ('Sharon Duncan-Brewster', 'female', '1976-02-08');

INSERT INTO movies (title, description, release_date, rating)
VALUES ('Oppenheimer',
        'The story of J. Robert Oppenheimer''s role in the development of the atomic bomb during World War II.',
        '2023-07-21', 8),
       ('Retreat',
        'Kate and Martin escape from personal tragedy to an Island Retreat. Cut off from the outside world, their attempts to recover are shattered when a man is washed ashore, with news of airborne killer disease that is sweeping through Europe.',
        '2011-10-14', 5),
       ('Dune: Part Two',
        'Follow the mythic journey of Paul Atreides as he unites with Chani and the Fremen while on a path of revenge against the conspirators who destroyed his family. Facing a choice between the love of his life and the fate of the known universe, Paul endeavors to prevent a terrible future only he can foresee.',
        '2024-02-29', 8),
       ('Dune',
        'Paul Atreides, a brilliant and gifted young man born into a great destiny beyond his understanding, must travel to the most dangerous planet in the universe to ensure the future of his family and his people. As malevolent forces explode into conflict over the planet''s exclusive supply of the most precious resource in existence-a commodity capable of unlocking humanity''s greatest potential-only those who can conquer their fear will survive.',
        '2021-10-01', 8),
       ('Space Jam: A New Legacy ',
        'When LeBron and his young son Dom are trapped in a digital space by a rogue A.I., LeBron must get them home safe by leading Bugs, Lola Bunny and the whole gang of notoriously undisciplined Looney Tunes to victory over the A.I.''s digitized champions on the court. It''s Tunes versus Goons in the highest-stakes challenge of his life.',
        '2021-07-16', 7);
INSERT INTO movies_actors (movie_id, actor_id)
VALUES (1, 1),
       (1, 3),
       (1, 2),
       (2, 1),
       (3, 4),
       (3, 5),
       (3, 6),
       (4, 5),
       (4, 6),
       (4, 7),
       (5, 5);