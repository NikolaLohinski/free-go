package client

import (
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

const (
	interfaceNotFoundCode     = "nodev"
	interfaceHostNotFoundCode = "nohost"
)

func (c *client) ListLanInterfaceInfo() (result []types.LanInfo, err error) {
	response, err := c.get("lan/browser/interfaces/", c.withSession)
	if err != nil {
		return nil, fmt.Errorf("failed to GET lan/browser/interfaces/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get lan info from generic response: %w", err)
	}

	return result, nil
}

func (c *client) GetLanInterface(name string) (result []types.LanInterfaceHost, err error) {
	response, err := c.get(fmt.Sprintf("lan/browser/%s", name), c.withSession)
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

func (c *client) GetLanInterfaceHost(interfaceName, identifier string) (result types.LanInterfaceHost, err error) {
	response, err := c.get(fmt.Sprintf("lan/browser/%s/%s", interfaceName, identifier), c.withSession)
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
