# Role System Documentation

## Overview

This SaaS template implements a **dual-role system** for flexible access control across two levels:

1. **System Roles** (Platform-level) - Controls platform administration
2. **Organization Roles** (Tenant-level) - Controls organization membership

## Design Philosophy

### Why VARCHAR Instead of ENUM?

**✅ Advantages:**
- **Flexibility**: Add new roles without database migration
- **Zero-downtime**: Change roles in code only
- **Type-safe**: Constants in Go provide compile-time safety
- **Extensible**: Easy to customize per use case

**Database Schema:**
```sql
-- Flexible VARCHAR with helpful comments
system_role VARCHAR(50) NOT NULL DEFAULT 'user'
role VARCHAR(50) NOT NULL DEFAULT 'member'
```

**Application Layer:**
```go
// Type-safe constants in Go
const (
    SystemRoleUser = "user"
    OrgRoleOwner   = "owner"
)
```

## System Roles (Platform Level)

Platform-level roles for SaaS owners to manage all tenants.

### Constants Location
`internal/entity/user_entity.go`

### Available Roles

| Role | Constant | Description | Use Case |
|------|----------|-------------|----------|
| **user** | `SystemRoleUser` | Regular tenant (default) | All customers |
| **support** | `SystemRoleSupport` | Customer support access | Support team |
| **admin** | `SystemRoleAdmin` | Platform administrator | Operations team |
| **super_admin** | `SystemRoleSuperAdmin` | Platform owner | Founders, CTO |

### Database Schema

```sql
-- users table
system_role VARCHAR(50) NOT NULL DEFAULT 'user'
INDEX idx_users_system_role (system_role)
```

### Usage Example

```go
// Check system role
if user.IsSuperAdmin() {
    // Full platform access
}

if user.IsSystemAdmin() {
    // Platform admin access
}

if user.IsSupport() {
    // Support access (can impersonate)
}

// Validate system role
if err := entity.ValidateSystemRole("admin"); err != nil {
    return err // Invalid role
}
```

### Permissions Matrix

| Action | User | Support | Admin | Super Admin |
|--------|------|---------|-------|-------------|
| Access own org | ✅ | ✅ | ✅ | ✅ |
| View all orgs | ❌ | ❌ | ✅ | ✅ |
| Impersonate users | ❌ | ✅ | ❌ | ✅ |
| Suspend orgs | ❌ | ❌ | ✅ | ✅ |
| Delete orgs | ❌ | ❌ | ❌ | ✅ |
| Platform settings | ❌ | ❌ | ❌ | ✅ |

## Organization Roles (Tenant Level)

Organization-level roles for managing team members within each tenant.

### Constants Location
`internal/entity/organization_member_entity.go`

### Available Roles

| Role | Constant | Description | Use Case |
|------|----------|-------------|----------|
| **owner** | `OrgRoleOwner` | Organization owner (1 per org) | Account owner, pays bills |
| **admin** | `OrgRoleAdmin` | Organization administrator | Team managers |
| **member** | `OrgRoleMember` | Regular member (default) | Regular employees |

### Database Schema

```sql
-- organization_members table
role VARCHAR(50) NOT NULL DEFAULT 'member'
INDEX idx_member_role (role)
```

### Usage Example

```go
// Check organization role
if member.IsOwner() {
    // Full organization control
}

if member.IsAdmin() {
    // Admin or owner access
}

// Validate organization role
if err := entity.ValidateOrganizationRole("admin"); err != nil {
    return err // Invalid role
}
```

### Permissions Matrix

| Action | Owner | Admin | Member |
|--------|-------|-------|--------|
| View org details | ✅ | ✅ | ✅ |
| Update org settings | ✅ | ✅ | ❌ |
| Manage subscription | ✅ | ❌ | ❌ |
| Invite members | ✅ | ✅ | ❌ |
| Remove members | ✅ | ✅ | ❌ |
| Change roles | ✅ | ✅* | ❌ |
| Delete organization | ✅ | ❌ | ❌ |
| Leave organization | ❌** | ✅ | ✅ |

**Notes:**
- `*` Admin can only change member ↔ admin, not owner
- `**` Owner must transfer ownership before leaving

## Extending Roles

### Adding New System Roles

**1. Update constants:**
```go
// internal/entity/user_entity.go
const (
    SystemRoleUser       = "user"
    SystemRoleSupport    = "support"
    SystemRoleAdmin      = "admin"
    SystemRoleSuperAdmin = "super_admin"
    SystemRoleBilling    = "billing"  // NEW ROLE
)
```

**2. Update validator:**
```go
// internal/entity/role_validator.go
func ValidSystemRoles() []string {
    return []string{
        SystemRoleUser,
        SystemRoleSupport,
        SystemRoleAdmin,
        SystemRoleSuperAdmin,
        SystemRoleBilling, // ADD HERE
    }
}
```

**3. No database migration needed!** ✅

### Adding New Organization Roles

**Example: Adding "viewer" role for read-only access**

```go
// internal/entity/organization_member_entity.go
const (
    OrgRoleOwner  = "owner"
    OrgRoleAdmin  = "admin"
    OrgRoleMember = "member"
    OrgRoleViewer = "viewer"  // NEW ROLE
)

// Add helper method
func (o *OrganizationMember) IsViewer() bool {
    return o.Role == OrgRoleViewer
}
```

**Update validator:**
```go
// internal/entity/role_validator.go
func ValidOrganizationRoles() []string {
    return []string{
        OrgRoleOwner,
        OrgRoleAdmin,
        OrgRoleMember,
        OrgRoleViewer, // ADD HERE
    }
}
```

## Use Case Examples

### Example 1: SaaS POS System

```go
// Custom labels for POS context
const (
    // Use generic constants
    RoleOwner  = entity.OrgRoleOwner  // Store Owner
    RoleManager = entity.OrgRoleAdmin  // Store Manager
    RoleCashier = entity.OrgRoleMember // Cashier
)

// Business logic
func CanProcessRefund(member *entity.OrganizationMember) bool {
    // Only owner and manager can refund
    return member.IsAdmin()
}

func CanProcessSale(member *entity.OrganizationMember) bool {
    // Everyone can process sales
    return member.IsMember()
}
```

### Example 2: Blog/CMS System

```go
// Custom labels for Blog context
const (
    RoleOwner  = entity.OrgRoleOwner  // Blog Owner
    RoleEditor = entity.OrgRoleAdmin  // Chief Editor
    RoleWriter = entity.OrgRoleMember // Content Writer
)

// Business logic
func CanPublishPost(member *entity.OrganizationMember) bool {
    // Only owner and editor can publish
    return member.IsAdmin()
}

func CanCreateDraft(member *entity.OrganizationMember) bool {
    // Everyone can create drafts
    return member.IsMember()
}
```

### Example 3: Project Management

```go
// Custom labels for PM context
const (
    RoleOwner   = entity.OrgRoleOwner  // Company Owner
    RoleManager = entity.OrgRoleAdmin  // Project Manager
    RoleDev     = entity.OrgRoleMember // Developer
)

// Business logic
func CanAssignTasks(member *entity.OrganizationMember) bool {
    // Only owner and PM can assign
    return member.IsAdmin()
}

func CanViewProjects(member *entity.OrganizationMember) bool {
    // Everyone can view
    return true
}
```

## Validation

### Server-Side Validation

```go
// Validate before creating organization member
func (u *UseCase) AddMember(orgID, userID, role string) error {
    // Validate role
    if err := entity.ValidateOrganizationRole(role); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, err.Error())
    }
    
    // Ensure only one owner
    if role == entity.OrgRoleOwner {
        exists, _ := u.Repository.HasOwner(orgID)
        if exists {
            return fiber.NewError(fiber.StatusConflict, "Organization already has an owner")
        }
    }
    
    // Create member
    // ...
}
```

### Request Validation

```go
// internal/model/organization_model.go
type AddMemberRequest struct {
    UserID string `json:"user_id" validate:"required,uuid"`
    Role   string `json:"role" validate:"required,oneof=owner admin member"`
}
```

## Security Best Practices

### 1. Always Validate Roles
```go
// ❌ BAD - Direct string assignment
member.Role = "admin"

// ✅ GOOD - Validate first
if err := entity.ValidateOrganizationRole(role); err != nil {
    return err
}
member.Role = role
```

### 2. Use Constants
```go
// ❌ BAD - Magic strings
if member.Role == "owner" {
    // ...
}

// ✅ GOOD - Constants
if member.Role == entity.OrgRoleOwner {
    // ...
}
```

### 3. Create Helper Methods
```go
// ✅ Encapsulate permission logic
if member.IsAdmin() {
    // Clear and maintainable
}
```

### 4. Audit Role Changes
```go
// Log important role changes
func (u *UseCase) ChangeRole(orgID, userID, newRole string) error {
    oldRole := member.Role
    member.Role = newRole
    
    // Audit log
    u.AuditLog.Create(&entity.AuditLog{
        Action: "role_changed",
        UserID: currentUser.ID,
        TargetUserID: userID,
        OldValue: oldRole,
        NewValue: newRole,
    })
    
    return nil
}
```

## Migration Guide

### Upgrading from ENUM to VARCHAR

If you previously used ENUM and want to switch:

```sql
-- Step 1: Add new VARCHAR column
ALTER TABLE organization_members 
ADD COLUMN role_new VARCHAR(50) NOT NULL DEFAULT 'member';

-- Step 2: Copy data
UPDATE organization_members SET role_new = role;

-- Step 3: Drop old ENUM column
ALTER TABLE organization_members DROP COLUMN role;

-- Step 4: Rename new column
ALTER TABLE organization_members 
CHANGE role_new role VARCHAR(50) NOT NULL DEFAULT 'member';

-- Step 5: Add index
ALTER TABLE organization_members ADD INDEX idx_member_role (role);
```

## Performance Considerations

### Indexing
```sql
-- Both role columns are indexed for fast filtering
INDEX idx_users_system_role (system_role)
INDEX idx_member_role (role)
```

### Query Performance
```sql
-- Fast with index
SELECT * FROM users WHERE system_role = 'admin';

-- Fast with index
SELECT * FROM organization_members 
WHERE organization_id = ? AND role = 'owner';
```

### Caching Strategy
```go
// Cache user with roles
type CachedUser struct {
    User       *entity.User
    OrgRole    string
    SystemRole string
    ExpiresAt  time.Time
}

// Cache key: user_id:org_id
func GetCachedUserRole(userID, orgID string) (*CachedUser, error) {
    key := fmt.Sprintf("user_role:%s:%s", userID, orgID)
    // ... Redis get
}
```

## FAQ

**Q: Why not use a separate roles table?**
A: For 3-5 fixed roles, VARCHAR is simpler and faster. Use a roles table only if you need:
- Custom roles per organization
- Hundreds of roles
- Dynamic permissions
- Role metadata (description, color, etc)

**Q: Can I rename roles without migration?**
A: Yes! Update the constants in code, then run a data migration:
```sql
UPDATE organization_members SET role = 'manager' WHERE role = 'admin';
```

**Q: How do I add role-specific data?**
A: Add optional JSON column:
```sql
role_metadata JSON NULL
-- Example: {"display_name": "Store Manager", "color": "#FF5733"}
```

**Q: Should I use ENUM for stricter validation?**
A: VARCHAR + application validation is more flexible. If you need strict DB-level validation, use CHECK constraint:
```sql
CHECK (role IN ('owner', 'admin', 'member'))
```

## Related Files

- `internal/entity/user_entity.go` - System role constants & helpers
- `internal/entity/organization_member_entity.go` - Org role constants & helpers
- `internal/entity/role_validator.go` - Role validation functions
- `db/migrations/000002_create_table_users.up.sql` - Users table with system_role
- `db/migrations/000003_create_table_organization_members.up.sql` - Members table with role
