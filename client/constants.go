package client

import (
	"time"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// Errors.
	ErrAppIDIsNotSet              = Error("app id is not set")
	ErrPrivateTokenIsNotSet       = Error("private token is not set")
	ErrInterfaceNotFound          = Error("interface not found")
	ErrInterfaceHostNotFound      = Error("interface host not found")
	ErrPortForwardingRuleNotFound = Error("port forwarding rule not found")
)

var (
	// Login.
	LoginSessionTTL = time.Minute * 30 // Fixed by the freebox server, but made into a variable for unit testing

	// Authorize.
	AuthorizeGrantingTimeout = time.Minute * 5
	AuthorizeRetryDelay      = time.Second * 5
)
