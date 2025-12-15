CREATE TABLE plans (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    billing_period VARCHAR(20) NOT NULL,
    features JSON,
    limits JSON,
    is_active BOOLEAN DEFAULT TRUE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT NULL DEFAULT NULL,
    INDEX idx_plan_slug (slug),
    INDEX idx_plan_active (is_active),
    INDEX idx_plan_deleted (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
