package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

const (
	interfaceNotFoundCode     = "nodev"
	interfaceHostNotFoundCode = "nohost"
)

func (c *client) ListLanInterfaceInfo(ctx context.Context) (result []types.LanInfo, err error) {
	response, err := c.get(ctx, "lan/browser/interfaces/", c.withSession(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to GET lan/browser/interfaces/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get lan info from generic response: %w", err)
	}

	return result, nil
}

func (c *client) GetLanInterface(ctx context.Context, name string) (result []types.LanInterfaceHost, err error) {
	response, err := c.get(ctx, fmt.Sprintf("lan/browser/%s", name), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == interfaceNotFoundCode {
			return result, ErrInterfaceNotFound
		}

		return result, fmt.Errorf("failed to GET lan/browser/%s endpoint: %w", name, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get lan interface from generic response: %w", err)
	}

	return result, nil
}

func (c *client) GetLanInterfaceHost(ctx context.Context, interfaceName, identifier string) (result types.LanInterfaceHost, err error) {
	response, err := c.get(ctx, fmt.Sprintf("lan/browser/%s/%s", interfaceName, identifier), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == interfaceNotFoundCode {
			return result, ErrInterfaceNotFound
		}

		if response != nil && response.ErrorCode == interfaceHostNotFoundCode {
			return result, ErrInterfaceHostNotFound
		}

		return result, fmt.Errorf("failed to GET lan/browser/%s/%s endpoint: %w", interfaceName, identifier, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get port lan interface host from generic response: %w", err)
	}

	return result, nil
}
