-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects(
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS teams(
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255),
    project_id UUID UNIQUE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users_in_teams(
    user_id UUID NOT NULL,
    team_id UUID NOT NULL,
    CONSTRAINT users_in_team_pkey PRIMARY KEY (user_id, team_id)
);

CREATE TABLE IF NOT EXISTS columns(
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255),
    project_id UUID NOT NULL,
    FOREIGN KEY(project_id) REFERENCES projects(id)
);

CREATE TABLE IF NOT EXISTS tasks(
    id UUID NOT NULL PRIMARY KEY,
    num INT NOT NULL,
    description VARCHAR(1000) NOT NULL,
    project_id UUID NOT NULL,
    executor_id UUID,
    column_id UUID,
    FOREIGN KEY(column_id) REFERENCES columns(id),
    FOREIGN KEY(project_id) REFERENCES projects(id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS columns;
DROP TABLE IF EXISTS users_in_teams;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS projects;
-- +goose StatementEnd
