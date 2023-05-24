package types

type LoginResponse struct {
	Message   string      `json:"msg,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
	Success   bool        `json:"success"`
	Result    LoginResult `json:"result"`
}

type LoginResult struct {
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
	UID       string        `json:"uid,omitempty"`
	Message   string        `json:"msg,omitempty"`
	ErrorCode string        `json:"error_code,omitempty"`
	Success   bool          `json:"success"`
	Result    SessionResult `json:"result"`
}

type SessionResult struct {
	SessionToken string          `json:"session_token,omitempty"`
	PasswordSet  bool            `json:"password_set,omitempty"`
	Permissions  map[string]bool `json:"permissions,omitempty"`
	Challenge    string          `json:"challenge"`
	PasswordSalt string          `json:"password_salt"`
}

type Permissions = map[string]bool
