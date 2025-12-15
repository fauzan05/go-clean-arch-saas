CREATE TABLE audit_logs (
    id CHAR(36) NOT NULL PRIMARY KEY,
    user_id CHAR(36),
    organization_id CHAR(36),
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    resource_id CHAR(36),
    details JSON,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at BIGINT NOT NULL,
    deleted_at BIGINT NULL DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL,
    INDEX idx_audit_user (user_id),
    INDEX idx_audit_org (organization_id),
    INDEX idx_audit_created (created_at),
    INDEX idx_audit_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
