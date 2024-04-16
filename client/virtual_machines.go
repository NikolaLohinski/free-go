package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

const (
	codeVirtualMachineNotFound = "no_such_vm"
)

func (c *client) GetVirtualMachineInfo(ctx context.Context) (result types.VirtualMachinesInfo, err error) {
	response, err := c.get(ctx, "vm/info/", c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to GET vm/info/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get vm info from generic response: %w", err)
	}

	return result, nil
}

func (c *client) GetVirtualMachineDistributions(ctx context.Context) (result []types.VirtualMachineDistribution, err error) {
	response, err := c.get(ctx, "vm/distros/", c.withSession(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to GET vm/distros/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get vm distributions from generic response: %w", err)
	}

	return result, nil
}

func (c *client) ListVirtualMachines(ctx context.Context) (result []types.VirtualMachine, err error) {
	response, err := c.get(ctx, "vm/", c.withSession(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to GET vm/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get vms from generic response: %w", err)
	}

	return result, nil
}

func (c *client) CreateVirtualMachine(ctx context.Context, payload types.VirtualMachinePayload) (result types.VirtualMachine, err error) {
	response, err := c.post(ctx, "vm/", payload, c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to POST to vm/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get vm from generic response: %w", err)
	}

	return result, nil
}

func (c *client) UpdateVirtualMachine(ctx context.Context, identifier int64, payload types.VirtualMachinePayload) (result types.VirtualMachine, err error) {
	response, err := c.put(ctx, fmt.Sprintf("vm/%d", identifier), payload, c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to PUT to vm/%d endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get vm from generic response: %w", err)
	}

	return result, nil
}

func (c *client) GetVirtualMachine(ctx context.Context, identifier int64) (result types.VirtualMachine, err error) {
	response, err := c.get(ctx, fmt.Sprintf("vm/%d", identifier), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeVirtualMachineNotFound {
			return result, ErrVirtualMachineNotFound
		}
		return result, fmt.Errorf("failed to GET to vm/%d endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get vm from generic response: %w", err)
	}

	return result, nil
}

func (c *client) DeleteVirtualMachine(ctx context.Context, identifier int64) error {
	response, err := c.delete(ctx, fmt.Sprintf("vm/%d", identifier), c.withSession(ctx))
	if err != nil {
		if response != nil && response.ErrorCode == codeVirtualMachineNotFound {
			return ErrVirtualMachineNotFound
		}
		return fmt.Errorf("failed to DELETE to vm/%d endpoint: %w", identifier, err)
	}

	return nil
}

func (c *client) StartVirtualMachine(ctx context.Context, identifier int64) error {
	if response, err := c.post(ctx, fmt.Sprintf("vm/%d/start", identifier), nil, c.withSession(ctx)); err != nil {
		if response != nil && response.ErrorCode == codeVirtualMachineNotFound {
			return ErrVirtualMachineNotFound
		}
		return fmt.Errorf("failed to POST to vm/%d/start endpoint: %w", identifier, err)
	}
	return nil
}

func (c *client) StopVirtualMachine(ctx context.Context, identifier int64) error {
	if response, err := c.post(ctx, fmt.Sprintf("vm/%d/stop", identifier), nil, c.withSession(ctx)); err != nil {
		if response != nil && response.ErrorCode == codeVirtualMachineNotFound {
			return ErrVirtualMachineNotFound
		}
		return fmt.Errorf("failed to POST to vm/%d/stop endpoint: %w", identifier, err)
	}
	return nil
}
