package entity

import "fmt"

// ValidSystemRoles returns all valid system roles
func ValidSystemRoles() []string {
	return []string{
		SystemRoleUser,
		SystemRoleSupport,
		SystemRoleAdmin,
		SystemRoleSuperAdmin,
	}
}

// ValidOrganizationRoles returns all valid organization roles
func ValidOrganizationRoles() []string {
	return []string{
		OrgRoleOwner,
		OrgRoleAdmin,
		OrgRoleMember,
	}
}

// IsValidSystemRole checks if the given system role is valid
func IsValidSystemRole(role string) bool {
	for _, validRole := range ValidSystemRoles() {
		if role == validRole {
			return true
		}
	}
	return false
}

// IsValidOrganizationRole checks if the given organization role is valid
func IsValidOrganizationRole(role string) bool {
	for _, validRole := range ValidOrganizationRoles() {
		if role == validRole {
			return true
		}
	}
	return false
}

// ValidateSystemRole validates system role and returns error if invalid
func ValidateSystemRole(role string) error {
	if !IsValidSystemRole(role) {
		return fmt.Errorf("invalid system role: %s, must be one of: %v", role, ValidSystemRoles())
	}
	return nil
}

// ValidateOrganizationRole validates organization role and returns error if invalid
func ValidateOrganizationRole(role string) error {
	if !IsValidOrganizationRole(role) {
		return fmt.Errorf("invalid organization role: %s, must be one of: %v", role, ValidOrganizationRoles())
	}
	return nil
}
