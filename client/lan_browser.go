package client

import (
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

func (c *client) ListLanInterfaceInfo() (result []types.LanInfo, err error) {
	response, err := c.Get("lan/browser/interfaces/", c.withSession)
	if err != nil {
		return nil, fmt.Errorf("failed to GET lan/browser/interfaces/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get port forwarding rules from generic response: %w", err)
	}

	return result, nil
}
