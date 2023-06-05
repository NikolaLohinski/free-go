package types

type LoginResponse struct {
	LoggedIn     bool   `json:"logged_in"`
	Challenge    string `json:"challenge"`
	PasswordSalt string `json:"password_salt"`
	PasswordSet  bool   `json:"password_set"`
}

type SessionsRequest struct {
	AppID    string `json:"app_id"`
	Password string `json:"password"`
}

type SessionResponse struct {
	SessionToken string          `json:"session_token,omitempty"`
	PasswordSet  bool            `json:"password_set,omitempty"`
	Permissions  map[string]bool `json:"permissions,omitempty"`
	Challenge    string          `json:"challenge"`
	PasswordSalt string          `json:"password_salt"`
}

type Permissions = map[string]bool
