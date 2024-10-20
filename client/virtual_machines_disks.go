package client

import (
	"context"
	"fmt"

	"github.com/nikolalohinski/free-go/types"
)

const (
	// Errors.
	ErrVMDiskSizeInvalid = Error("vm disk size is invalid")
)

// GetVirtualDiskInfo gets a disk info.
func (c *client) GetVirtualDiskInfo(ctx context.Context, path string) (result types.VirtualDiskInfo, err error) {
	response, err := c.post(ctx, "vm/disk/info/", &types.GetVirtualDiskPayload{
		DiskPath: types.Base64Path(path),
	}, c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to GET vm/disk/info endpoint: %w", err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get a vm disk info from generic response: %w", err)
	}

	return result, nil
}

// CreateVirtualDisk creates a new disk.
func (c *client) CreateVirtualDisk(ctx context.Context, payload types.VirtualDisksCreatePayload) (result int64, err error) {
	if payload.Size < 0 {
		return result, ErrVMDiskSizeInvalid
	}

	response, err := c.post(ctx, "vm/disk/create/", payload, c.withSession(ctx))
	if err != nil {
		return 0, fmt.Errorf("failed to GET vm/disk/create endpoint: %w", err)
	}

	var responseBody struct {
		ID int64 `json:"id"`
	}
	if err = c.fromGenericResponse(response, &responseBody); err != nil {
		return 0, fmt.Errorf("failed to get an ID from generic response: %w", err)
	}

	return responseBody.ID, nil
}

// GetVirtualDiskTask gets a disk task.
func (c *client) GetVirtualDiskTask(ctx context.Context, identifier int64) (result types.VirtualMachineDiskTask, err error) {
	response, err := c.get(ctx, fmt.Sprintf("vm/disk/task/%d", identifier), c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to GET vm/disk/task/%d endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get a disk task from generic response: %w", err)
	}

	return result, nil
}

// ResizeVirtualDisk resizes a existing disk.
func (c *client) ResizeVirtualDisk(ctx context.Context, payload types.VirtualDisksResizePayload) (result int64, err error) {
	if payload.NewSize < 0 {
		return result, ErrVMDiskSizeInvalid
	}

	response, err := c.post(ctx, "vm/disk/resize/", payload, c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to GET vm/disk/resize endpoint: %w", err)
	}

	var responseBody struct {
		ID int64 `json:"id"`
	}
	if err = c.fromGenericResponse(response, &responseBody); err != nil {
		return 0, fmt.Errorf("failed to get an ID from generic response: %w", err)
	}

	return responseBody.ID, nil
}

// DeleteVirtualDiskTask deletes a disk task once done.
func (c *client) DeleteVirtualDiskTask(ctx context.Context, identifier int64) error {
	_, err := c.delete(ctx, fmt.Sprintf("vm/disk/task/%d", identifier), c.withSession(ctx))
	if err != nil {
		return fmt.Errorf("failed to GET vm/disk/task/%d endpoint: %w", identifier, err)
	}

	return nil
}

// GetVirtualMachineDiskTask gets a disk task.
func (c *client) GetVirtualMachineDiskTask(ctx context.Context, identifier int64) (result types.VirtualMachineDiskTask, err error) {
	response, err := c.get(ctx, fmt.Sprintf("vm/disk/task/%d", identifier), c.withSession(ctx))
	if err != nil {
		return result, fmt.Errorf("failed to GET vm/disk/task/%d endpoint: %w", identifier, err)
	}

	if err = c.fromGenericResponse(response, &result); err != nil {
		return result, fmt.Errorf("failed to get a disk task from generic response: %w", err)
	}

	return result, nil
}
