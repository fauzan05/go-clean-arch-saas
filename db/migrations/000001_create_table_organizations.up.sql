CREATE TABLE organizations (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT NULL DEFAULT NULL,
    INDEX idx_org_slug (slug),
    INDEX idx_org_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
