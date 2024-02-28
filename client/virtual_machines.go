package client

import (
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

const (
	codeVirtualMachineNotFound = "no_such_vm"
)

func (c *client) GetVirtualMachineInfo() (result types.VirtualMachinesInfo, err error) {
	response, err := c.get("vm/info/", c.withSession)
	if err != nil {
		return result, fmt.Errorf("failed to GET vm/info/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get vm info from generic response: %w", err)
	}

	return result, nil
}

func (c *client) GetVirtualMachineDistributions() (result []types.VirtualMachineDistribution, err error) {
	response, err := c.get("vm/distros/", c.withSession)
	if err != nil {
		return nil, fmt.Errorf("failed to GET vm/distros/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get vm distributions from generic response: %w", err)
	}

	return result, nil
}

func (c *client) ListVirtualMachines() (result []types.VirtualMachine, err error) {
	response, err := c.get("vm/", c.withSession)
	if err != nil {
		return nil, fmt.Errorf("failed to GET vm/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return nil, fmt.Errorf("failed to get vms from generic response: %w", err)
	}

	return result, nil
}

func (c *client) CreateVirtualMachine(payload types.VirtualMachinePayload) (result types.VirtualMachine, err error) {
	response, err := c.post("vm/", payload, c.withSession)
	if err != nil {
		return result, fmt.Errorf("failed to POST to vm/ endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get vm from generic response: %w", err)
	}

	return result, nil
}

func (c *client) GetVirtualMachine(identifier int64) (result types.VirtualMachine, err error) {
	response, err := c.get(fmt.Sprintf("vm/%d", identifier), c.withSession)
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

func (c *client) DeleteVirtualMachine(identifier int64) error {
	response, err := c.delete(fmt.Sprintf("vm/%d", identifier), c.withSession)
	if err != nil {
		if response != nil && response.ErrorCode == codeVirtualMachineNotFound {
			return ErrVirtualMachineNotFound
		}
		return fmt.Errorf("failed to DELETE to vm/%d endpoint: %w", identifier, err)
	}

	return nil
}
