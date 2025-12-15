-- Seed default plans
INSERT INTO plans (id, name, slug, price, billing_period, features, limits, created_at, updated_at) VALUES
('plan-free', 'Free', 'free', 0.00, 'monthly', 
 '{"storage": "1GB", "users": "1", "support": "Community"}',
 '{"api_calls_per_month": 1000, "max_users": 1, "storage_gb": 1}',
 UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
 
('plan-pro', 'Pro', 'pro', 29.00, 'monthly',
 '{"storage": "50GB", "users": "10", "support": "Email"}',
 '{"api_calls_per_month": 100000, "max_users": 10, "storage_gb": 50}',
 UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000),
 
('plan-enterprise', 'Enterprise', 'enterprise', 99.00, 'monthly',
 '{"storage": "Unlimited", "users": "Unlimited", "support": "Priority"}',
 '{"api_calls_per_month": -1, "max_users": -1, "storage_gb": -1}',
 UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000);

-- Sample organization
INSERT INTO organizations (id, name, slug, created_at, updated_at) VALUES
('org-demo', 'Demo Organization', 'demo-org', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000);

-- Sample user (password: password123)
INSERT INTO users (id, name, email, password, organization_id, created_at, updated_at) VALUES
('user-demo', 'Demo User', 'demo@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'org-demo', UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000);

-- Assign user to org
INSERT INTO organization_members (organization_id, user_id, role, joined_at) VALUES
('org-demo', 'user-demo', 'owner', UNIX_TIMESTAMP() * 1000);

-- Free subscription for demo org
INSERT INTO subscriptions (id, organization_id, plan_id, status, current_period_start, current_period_end, created_at, updated_at) VALUES
('sub-demo', 'org-demo', 'plan-free', 'active', UNIX_TIMESTAMP() * 1000, (UNIX_TIMESTAMP() + 2592000) * 1000, UNIX_TIMESTAMP() * 1000, UNIX_TIMESTAMP() * 1000);
