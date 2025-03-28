CREATE TABLE teams
(
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(50)  NOT NULL,
    cluster_name VARCHAR(100) NOT NULL,
    namespace    VARCHAR(100) NOT NULL
);

CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE                            NOT NULL,
    name     VARCHAR(50)                                   NOT NULL,
    surname  VARCHAR(50)                                   NOT NULL,
    team_id  INTEGER REFERENCES teams (id) ON DELETE CASCADE,
    role     VARCHAR(50) CHECK (role IN ('admin', 'user')) NOT NULL,
    password VARCHAR(255)
);

CREATE TABLE configs
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR,
    team_id    INTEGER REFERENCES teams (id) ON DELETE CASCADE,
    type       VARCHAR(50) CHECK (type IN ('ConfigMap', 'Secret')),
    content    TEXT                                                NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE config_versions
(
    id               SERIAL PRIMARY KEY,
    actual_config_id INTEGER REFERENCES configs (id) ON DELETE CASCADE
);

CREATE TABLE config_changes
(
    id         SERIAL PRIMARY KEY,
    new_config INTEGER                                              REFERENCES configs (id) ON DELETE SET NULL,
    old_config INTEGER                                              REFERENCES config_changes (id) ON DELETE SET NULL,
    user_id    INTEGER REFERENCES users (id) ON DELETE CASCADE,
    team_id    INTEGER REFERENCES teams (id) ON DELETE CASCADE,
    action     VARCHAR(50) CHECK (action IN ('create', 'rollback')) NOT NULL,
    action_at  TIMESTAMP DEFAULT now()
);
