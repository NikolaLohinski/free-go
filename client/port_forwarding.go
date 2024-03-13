package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

const (
	codePortForwardingNotFound = "noent"
)

func (c *client) ListPortForwardingRules(ctx context.Context) ([]types.PortForwardingRule, error) {
	response, err := c.get(ctx, "fw/redir/", c.withSession(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to GET fw/redir/ endpoint: %w", err)
	}

	result := make([]types.PortForwardingRule, 0)
	if response.Result != nil {
		if err = c.fromGenericResponse(response, &result); err != nil {
			return nil, fmt.Errorf("failed to get port forwarding rules from generic response: %w", err)
		}
	}

	return result, nil
}

func (c *client) GetPortForwardingRule(ctx context.Context, identifier int64) (rule types.PortForwardingRule, err error) {
	response, err := c.get(ctx, fmt.Sprintf("fw/redir/%d", identifier), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codePortForwardingNotFound {
			return rule, ErrPortForwardingRuleNotFound
		}

		return rule, fmt.Errorf("failed to GET fw/redir/%d endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &rule); err != nil {
		return rule, fmt.Errorf("failed to get a port forwarding rule from a generic response: %w", err)
	}

	return rule, nil
}

func (c *client) CreatePortForwardingRule(
	ctx context.Context,
	payload types.PortForwardingRulePayload,
) (rule types.PortForwardingRule, err error) {
	response, err := c.post(ctx, "fw/redir/", payload, c.withSession(ctx))
	if err != nil {
		return rule, fmt.Errorf("failed to POST to fw/redir/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &rule); err != nil {
		return rule, fmt.Errorf("failed to get a port forwarding rule from a generic response: %w", err)
	}

	return rule, nil
}

func (c *client) DeletePortForwardingRule(ctx context.Context, identifier int64) error {
	response, err := c.delete(ctx, fmt.Sprintf("fw/redir/%d", identifier), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codePortForwardingNotFound {
			return ErrPortForwardingRuleNotFound
		}

		return fmt.Errorf("failed to GET fw/redir/%d endpoint: %w", identifier, err)
	}

	return nil
}

func (c *client) UpdatePortForwardingRule(
	ctx context.Context,
	identifier int64,
	payload types.PortForwardingRulePayload,
) (rule types.PortForwardingRule, err error) {
	response, err := c.put(ctx, fmt.Sprintf("fw/redir/%d", identifier), payload, c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codePortForwardingNotFound {
			return rule, ErrPortForwardingRuleNotFound
		}

		return rule, fmt.Errorf("failed to GET fw/redir/%d endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &rule); err != nil {
		return rule, fmt.Errorf("failed to get a port forwarding rule from a generic response: %w", err)
	}

	return rule, nil
}
