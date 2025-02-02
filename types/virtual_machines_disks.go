package types

type virtualMachineDiskTaskType string

const (
	DiskTaskTypeCreate virtualMachineDiskTaskType = "create"
	DiskTaskTypeResize virtualMachineDiskTaskType = "resize"
)

const (
	EventSourceVMDisk eventSource = "vm" // Disk events are sourced from the VM
	EventDiskTaskDone eventName   = "disk_task_done"
)

type diskError = string

const (
	DiskErrorNotFound     diskError = "file_not_found" // File not found
	DiskTaskErrorNotFound diskError = "task_notfound"  // Task not found
	DiskErrorInfo         diskError = "info"           // Unable to retrieve disk info
)

type diskTaskError string

const (
	// Undocumented and reverse engineered error codes.
	DiskTaskErrorNone     diskTaskError = "none"      // No error
	DiskTaskErrorTooSmall diskTaskError = "too_small" // Disk is too small
	DiskTaskErrorExists   diskTaskError = "exists"    // Disk already exists
	DiskTaskErrorInternal diskTaskError = "internal"  // Internal error (e.g. empty path on creation)
)

type VirtualDiskInfo struct {
	Type        diskType `json:"type"`
	ActualSize  int64    `json:"actual_size"`  // Space used by virtual image on disk. This is how much filesystem space is consumed on the box.
	VirtualSize int64    `json:"virtual_size"` // Size of virtual disk. This is the size the disk will appear inside the VM.
}

type VirtualDisksCreatePayload struct {
	DiskPath Base64Path `json:"disk_path"` // Base64 encoded
	Size     int64      `json:"size"`      // Size of virtual disk in bytes
	DiskType diskType   `json:"disk_type"`
}

type VirtualDisksResizePayload struct {
	DiskPath    Base64Path `json:"disk_path"`    // Base64 encoded
	NewSize     int64      `json:"size"`         // New size of virtual disk in bytes
	ShrinkAllow bool       `json:"shrink_allow"` // Whether shrinking the disk is allowed. Setting to true means this operation can be destructive.
}

type VirtualMachineDiskTask struct {
	ID    int64                      `json:"id"`
	Type  virtualMachineDiskTaskType `json:"type"`
	Done  bool                       `json:"done"`
	Error bool                       `json:"error"`
}

type GetVirtualDiskPayload struct {
	DiskPath Base64Path `json:"disk_path"` // Base64 encoded
}
