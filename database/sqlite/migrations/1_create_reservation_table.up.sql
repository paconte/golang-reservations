CREATE TABLE IF NOT EXISTS reservations (
    id INTEGER PRIMARY KEY NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    duration INT NOT NULL
);