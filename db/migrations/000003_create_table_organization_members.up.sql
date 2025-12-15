CREATE TABLE organization_members (
    organization_id VARCHAR(100) NOT NULL,
    user_id VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    joined_at BIGINT NOT NULL,
    PRIMARY KEY (organization_id, user_id),
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_member_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
