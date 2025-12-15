CREATE TABLE users (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    refresh_token VARCHAR(100) NULL,
    refresh_token_expires_at BIGINT NULL,
    organization_id VARCHAR(100) NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL,
    INDEX idx_users_email (email),
    INDEX idx_users_refresh_token (refresh_token),
    INDEX idx_users_org (organization_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
