CREATE TYPE submission_result AS ENUM (
    'compilation_error',
    'runtime_error',
    'wrong_answer',
    'time_limit_exceeded',
    'memory_limit_exceeded',
    'bad_code',
    'accepted'
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    user_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE problems (
    problem_id SERIAL PRIMARY KEY,
    author_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    constraints TEXT NOT NULL,
    time_limit INTEGER NOT NULL,
    memory_limit INTEGER NOT NULL
);

CREATE TABLE inputs (
    input_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    problem_id INTEGER NOT NULL REFERENCES problems(problem_id) ON DELETE CASCADE,
    input_storage_id INTEGER NOT NULL
);

CREATE TABLE submissions (
    submission_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    problem_id INTEGER NOT NULL REFERENCES problems(problem_id) ON DELETE CASCADE,
    output_id INTEGER,
    result submission_result NOT NULL,
    user_code_storage_id TEXT NOT NULL,
    time_stamp BIGINT NOT NULL
);

CREATE TABLE outputs (
    output_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    problem_id INTEGER NOT NULL REFERENCES problems(problem_id) ON DELETE CASCADE,
    output_storage_id INTEGER NOT NULL
);

CREATE TABLE checkers (
    checker_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    problem_id INTEGER NOT NULL REFERENCES problems(problem_id) ON DELETE CASCADE,
    checker_storage_id INTEGER NOT NULL
);

ALTER TABLE submissions
    ADD CONSTRAINT submissions_output_id_fkey
    FOREIGN KEY (output_id)
    REFERENCES outputs(output_id)
    ON DELETE SET NULL
    DEFERRABLE INITIALLY DEFERRED;

CREATE INDEX inputs_problem_id_idx ON inputs(problem_id);
CREATE INDEX inputs_user_id_idx ON inputs(user_id);
CREATE INDEX submissions_user_id_idx ON submissions(user_id);
CREATE INDEX submissions_problem_id_idx ON submissions(problem_id);
CREATE INDEX outputs_problem_id_idx ON outputs(problem_id);
CREATE INDEX outputs_user_id_idx ON outputs(user_id);
CREATE INDEX checkers_problem_id_idx ON checkers(problem_id);
CREATE INDEX checkers_user_id_idx ON checkers(user_id);
