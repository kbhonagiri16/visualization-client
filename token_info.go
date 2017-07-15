package openstack

import (
	"time"
)

const adminRoleName = "admin"

// TokenInfo represents data parsed from openstack token
type TokenInfo struct {
	ID          string
	ExpiresAt   time.Time
	Roles       []map[string]string
	ProjectID   string
	ProjectName string
}

// IsAdmin tries to find 'admin' role in Roles map
func (t *TokenInfo) IsAdmin() bool {
	for _, role := range t.Roles {
		if role["name"] == adminRoleName {
			return true
		}
	}
	return false
}
