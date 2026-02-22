-- +migrate Up
CREATE TABLE boards (
    id SERIAL PRIMARY KEY,
    external_id VARCHAR(36) NOT NULL UNIQUE,
    workspace_id INT NOT NULL,
    created_by_id INT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    active_status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255),
    CONSTRAINT fk_boards_workspace FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE,
    CONSTRAINT fk_boards_creator FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE SET NULL
);

-- +migrate Down
DROP TABLE boards;
