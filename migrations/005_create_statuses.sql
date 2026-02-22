-- +migrate Up
CREATE TABLE statuses (
    id SERIAL PRIMARY KEY,
    external_id VARCHAR(36) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    color VARCHAR(50),
    position INT DEFAULT 0,
    active_status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

-- +migrate Down
DROP TABLE statuses;
