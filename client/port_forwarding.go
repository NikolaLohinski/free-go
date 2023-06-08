package client

import (
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

func (c *client) PortForwardingRules() ([]types.PortForwardingRule, error) {
	response, err := c.genericGet("fw/redir/", c.withSession)
	if err != nil {
		return nil, fmt.Errorf("failed to GET fw/redir/ endpoint: %w", err)
	}

	result := make([]types.PortForwardingRule, 0)
	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get port forwarding rules from generic response: %w", err)
	}

	return result, nil
}
