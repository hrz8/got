-- +goose Up
-- +goose StatementBegin
CREATE TYPE ProjectEnv AS ENUM ('LIVE', 'SANDBOX');

CREATE TABLE IF NOT EXISTS projects (
    id SERIAL NOT NULL,
    public_id UUID NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(100),
    timezone VARCHAR(50) NOT NULL,
    environment ProjectEnv NOT NULL DEFAULT 'SANDBOX',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT projects_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX projects_public_id_key ON projects(public_id);

CREATE TABLE IF NOT EXISTS project_api_key (
    id SERIAL NOT NULL,
    public_id UUID NOT NULL,
    project_id INTEGER NOT NULL,
    name VARCHAR(50) NOT NULL,
    expiration TIMESTAMP NOT NULL,
    key VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT project_api_key_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX project_api_key_public_id_key ON project_api_key(public_id);

ALTER TABLE project_api_key ADD CONSTRAINT project_api_key_project_id_fkey FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE RESTRICT ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_api_key DROP CONSTRAINT IF EXISTS project_api_key_project_id_fkey;

DROP TABLE IF EXISTS project_api_key;

DROP INDEX IF EXISTS project_api_key_public_id_key;

DROP TABLE IF EXISTS projects;

DROP INDEX IF EXISTS projects_public_id_key;

DROP TYPE IF EXISTS ProjectEnv;
-- +goose StatementEnd
