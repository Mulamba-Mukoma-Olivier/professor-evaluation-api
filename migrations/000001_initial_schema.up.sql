-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    matricule TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL
);

-- Create professors table
CREATE TABLE IF NOT EXISTS professors (
    id BIGSERIAL PRIMARY KEY,
    matricule TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT,
    department TEXT NOT NULL,
    grade TEXT,
    active BOOLEAN,
    status TEXT DEFAULT 'active'
);

-- Create courses table
CREATE TABLE IF NOT EXISTS courses (
    id BIGSERIAL PRIMARY KEY,
    code TEXT UNIQUE,
    name TEXT,
    description TEXT,
    department TEXT,
    academic_year TEXT
);

-- Create criterions table
CREATE TABLE IF NOT EXISTS criterions (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    max_score BIGINT NOT NULL,
    active BOOLEAN
);

-- Create eligibilities table
CREATE TABLE IF NOT EXISTS eligibilities (
    id BIGSERIAL PRIMARY KEY,
    student_id BIGINT NOT NULL,
    eligible BOOLEAN,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create evaluations table
CREATE TABLE IF NOT EXISTS evaluations (
    id BIGSERIAL PRIMARY KEY,
    student_id BIGINT NOT NULL,
    professor_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    academic_year TEXT NOT NULL,
    period TEXT NOT NULL,
    submitted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (professor_id) REFERENCES professors(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Create evaluation_answers table
CREATE TABLE IF NOT EXISTS evaluation_answers (
    id BIGSERIAL PRIMARY KEY,
    evaluation_id BIGINT NOT NULL,
    criterion_id BIGINT NOT NULL,
    score BIGINT NOT NULL,
    FOREIGN KEY (evaluation_id) REFERENCES evaluations(id) ON DELETE CASCADE,
    FOREIGN KEY (criterion_id) REFERENCES criterions(id) ON DELETE CASCADE
);
