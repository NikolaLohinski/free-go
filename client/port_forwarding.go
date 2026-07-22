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

// createPortForwardingRulePayload adds the fields the Freebox OS web UI always
// sends when creating a fw/redir/ rule (id, hostname, host, valid) on top of the
// documented payload. Omitting them causes the API to reject the call with a
// generic "internal_error", even though those fields are read-only/computed on
// reads. They're all zero-valued placeholders since no host has been resolved yet.
type createPortForwardingRulePayload struct {
	types.PortForwardingRulePayload

	ID       int64  `json:"id"`
	Hostname string `json:"hostname"`
	Host     string `json:"host"`
	Valid    bool   `json:"valid"`
}

func (c *client) CreatePortForwardingRule(
	ctx context.Context,
	payload types.PortForwardingRulePayload,
) (rule types.PortForwardingRule, err error) {
	response, err := c.post(ctx, "fw/redir/", createPortForwardingRulePayload{PortForwardingRulePayload: payload}, c.withSession(ctx))
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

		return fmt.Errorf("failed to DELETE fw/redir/%d endpoint: %w", identifier, err)
	}

	return nil
}

// updatePortForwardingRulePayload adds the fields the Freebox OS web UI always
// sends when updating a fw/redir/ rule (id, hostname, host, valid) on top of the
// documented payload. Unlike on create, the API requires the rule's current host
// binding to be echoed back verbatim here: sending zero-valued placeholders is
// rejected with the same generic "internal_error" as omitting the fields entirely.
type updatePortForwardingRulePayload struct {
	types.PortForwardingRulePayload

	ID       int64                   `json:"id"`
	Hostname string                  `json:"hostname"`
	Host     *types.LanInterfaceHost `json:"host"`
	Valid    bool                    `json:"valid"`
}

func (c *client) UpdatePortForwardingRule(
	ctx context.Context,
	identifier int64,
	payload types.PortForwardingRulePayload,
) (rule types.PortForwardingRule, err error) {
	// The API rejects an update unless the current host binding (hostname, host,
	// valid) is echoed back verbatim alongside the changed fields, so the current
	// rule has to be read before it can be written back.
	current, err := c.GetPortForwardingRule(ctx, identifier)
	if err != nil {
		return rule, err
	}

	writePayload := updatePortForwardingRulePayload{
		PortForwardingRulePayload: payload,
		ID:                        identifier,
		Hostname:                  current.Hostname,
		Host:                      current.Host,
		Valid:                     current.Valid,
	}

	response, err := c.put(ctx, fmt.Sprintf("fw/redir/%d", identifier), writePayload, c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codePortForwardingNotFound {
			return rule, ErrPortForwardingRuleNotFound
		}

		return rule, fmt.Errorf("failed to PUT fw/redir/%d endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &rule); err != nil {
		return rule, fmt.Errorf("failed to get a port forwarding rule from a generic response: %w", err)
	}

	return rule, nil
}
