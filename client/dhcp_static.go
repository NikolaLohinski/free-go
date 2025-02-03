package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

func (c *client) ListDHCPStaticLease(ctx context.Context) (result []types.DHCPStaticLeaseInfo, err error) {
	response, err := c.get(ctx, "dhcp/static_lease/", c.withSession(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to GET dhcp/static_lease/ endpoint: %w", err)
	}

	if response.Result == nil {
		return
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to list dhcp static lease from generic response: %w", err)
	}

	return result, nil
}

func (c *client) GetDHCPStaticLease(ctx context.Context, identifier string) (result types.DHCPStaticLeaseInfo, err error) {
	response, err := c.get(ctx, "dhcp/static_lease/"+identifier, c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to GET dhcp/static_lease/%s endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get dhcp static lease from generic response: %w", err)
	}

	return result, nil
}

func (c *client) UpdateDHCPStaticLease(ctx context.Context, identifier string, payload types.DHCPStaticLeasePayload) (result types.LanInterfaceHost, err error) {
	response, err := c.put(ctx, "dhcp/static_lease/"+identifier, c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to PUT dhcp/static_lease/%s endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to update dhcp static lease from generic response: %w", err)
	}

	return result, nil
}

func (c *client) CreateDHCPStaticLease(ctx context.Context, payload types.DHCPStaticLeasePayload) (result types.LanInterfaceHost, err error) {
	response, err := c.post(ctx, "dhcp/static_lease/", payload, c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to POST dhcp/static_lease/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to create dhcp static lease from generic response: %w", err)
	}

	return result, nil
}

func (c *client) DeleteDHCPStaticLease(ctx context.Context, identifier string) error {
	_, err := c.delete(ctx, "dhcp/static_lease/"+identifier, c.withSession(ctx))
	if err != nil {
		return fmt.Errorf("failed to DELETE dhcp/static_lease/%s endpoint: %w", identifier, err)
	}

	return nil
}
