package types

type LoginRequest struct {
	UserName string `json:"user"`
	Password string `json:"password"`
}

type ClientRoles struct {
	Roles []string `json:"roles"`
}

type RealmAccess struct {
	ClientRoles
}

type Claims struct {
	RealmAccess       `json:"realm_access"`
	ResourceAccess    map[string]ClientRoles `json:"resource_access"`
	Groups            []string               `json:"cognito:groups"`
	Scope             string                 `json:"scope"`
	EmailVerified     bool                   `json:"email_verified"`
	UserName          string                 `json:"name"`
	PreferredUserName string                 `json:"preferred_username"`
	Email             string                 `json:"email"`
}
