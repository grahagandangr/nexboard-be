-- +migrate Up
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    external_id VARCHAR(36) NOT NULL UNIQUE,
    board_id INT NOT NULL,
    status_id INT NOT NULL,
    assigned_to INT,
    created_by_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    priority VARCHAR(50) NOT NULL,
    due_date TIMESTAMP,
    position INT DEFAULT 0,
    active_status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255),
    CONSTRAINT fk_tasks_board FOREIGN KEY (board_id) REFERENCES boards (id) ON DELETE CASCADE,
    CONSTRAINT fk_tasks_status FOREIGN KEY (status_id) REFERENCES statuses (id) ON DELETE RESTRICT,
    CONSTRAINT fk_tasks_assignee FOREIGN KEY (assigned_to) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT fk_tasks_creator FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE RESTRICT
);

-- +migrate Down
DROP TABLE tasks;
