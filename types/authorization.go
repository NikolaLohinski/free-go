package types

type PrivateToken = string

type AuthorizationRequest struct {
	Name    string
	Version string
	Device  string
}

type ErrorCode string

const (
	AuthorizationErrorCode ErrorCode = "auth_required" // "Vous devez vous connecter pour accéder à cette fonction"
)
