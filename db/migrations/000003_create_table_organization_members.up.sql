-- Valid role values: 'owner', 'admin', 'member'
-- These roles are for organization-level access control (not platform-level)
-- 'owner': Organization owner, full control including billing and deletion (1 per org)
-- 'admin': Organization administrator, can manage members and settings (but not billing)
-- 'member': Regular member, limited permissions (default)
-- Note: Extend this list as needed for your specific use case (e.g., 'viewer', 'manager')
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
    INDEX idx_member_role (role),
    INDEX idx_member_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
