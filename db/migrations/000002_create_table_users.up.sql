-- Valid system_role values: 'user', 'support', 'admin', 'super_admin'
-- These roles are for platform-level access control (not organization-level)
-- 'user' (default): Regular tenant user, only access their own organization
-- 'support': Customer support, can impersonate users and view organizations
-- 'admin': Platform admin, can view all organizations and system analytics
-- 'super_admin': Platform owner, full access to everything
CREATE TABLE users (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    system_role VARCHAR(50) NOT NULL DEFAULT 'user',
    email_verified TINYINT(1) DEFAULT 0,
    email_verified_at BIGINT NULL,
    verification_token VARCHAR(100) NULL,
    refresh_token VARCHAR(100) NULL,
    refresh_token_expires_at BIGINT NULL,
    organization_id CHAR(36) NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT NULL DEFAULT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL,
    INDEX idx_users_email (email),
    INDEX idx_users_system_role (system_role),
    INDEX idx_users_verification_token (verification_token),
    INDEX idx_users_refresh_token (refresh_token),
    INDEX idx_users_org (organization_id),
    INDEX idx_users_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
