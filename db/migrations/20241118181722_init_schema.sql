-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects(
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS users_in_projects(
    user_id UUID NOT NULL,
    project_id UUID NOT NULL,
    CONSTRAINT users_in_projects_pkey PRIMARY KEY (user_id, project_id), 
    FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_columns(
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255),
    project_id UUID NOT NULL,
    FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tasks(
    id UUID NOT NULL PRIMARY KEY,
    num INT NOT NULL,
    description VARCHAR(1000) NOT NULL,
    project_id UUID NOT NULL,
    executor_id UUID,
    column_id UUID,
    FOREIGN KEY(column_id) REFERENCES project_columns(id),
    FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS project_columns;
DROP TABLE IF EXISTS users_in_projects;
DROP TABLE IF EXISTS projects;
-- +goose StatementEnd
