package types

type Role struct {
	RoleName map[string]RoleConfig `json:"role"`
}

type RoleConfig struct {
	Path           string          `json:"path"`
	AllowedMethods map[string]bool `json:"allowed_methods"`
}
