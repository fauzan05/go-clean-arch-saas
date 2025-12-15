CREATE TABLE organizations (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_org_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
