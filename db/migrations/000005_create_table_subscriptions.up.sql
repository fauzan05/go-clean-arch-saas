CREATE TABLE subscriptions (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    organization_id VARCHAR(100) NOT NULL,
    plan_id VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    current_period_start BIGINT NOT NULL,
    current_period_end BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (plan_id) REFERENCES plans(id),
    INDEX idx_sub_org (organization_id),
    INDEX idx_sub_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
