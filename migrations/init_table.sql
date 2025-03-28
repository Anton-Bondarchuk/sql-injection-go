CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INT NOT NULL,
    sex BOOLEAN NOT NULL,
    card_id INT UNIQUE -- Один студент может иметь только одну карту
);

CREATE TABLE IF NOT EXISTS card_credits (
    id SERIAL PRIMARY KEY,
    student_id INT REFERENCES students(id) ON DELETE CASCADE,
    card_number BIGINT NOT NULL UNIQUE,
    expiration INT NOT NULL,
    cvv INT NOT NULL
);