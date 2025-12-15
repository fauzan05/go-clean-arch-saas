CREATE TABLE organization_members (
    organization_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    joined_at BIGINT NOT NULL,
    deleted_at BIGINT NULL DEFAULT NULL,
    PRIMARY KEY (organization_id, user_id),
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_member_user (user_id),
    INDEX idx_member_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
