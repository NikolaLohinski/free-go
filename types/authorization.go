package types

type PrivateToken = string

type AuthorizationRequest struct {
	Name    string
	Version string
	Device  string
}
