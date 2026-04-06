package types

type diskStatus string

const (
	DiskStatusNotDetected  diskStatus = "not_detected"
	DiskStatusDisabled     diskStatus = "disabled"
	DiskStatusInitializing diskStatus = "initializing"
	DiskStatusError        diskStatus = "error"
	DiskStatusActive       diskStatus = "active"
)

type boxFlavor string

const (
	BoxFlavorFull  boxFlavor = "full"
	BoxFlavorLight boxFlavor = "light"
)

type SystemConfig struct {
	FirmwareVersion  string     `json:"firmware_version"`  // freebox firmware version
	Mac              string     `json:"mac"`               // freebox mac address
	Serial           string     `json:"serial"`            // freebox serial number
	Uptime           string     `json:"uptime"`            // readable freebox uptime
	UptimeVal        int64      `json:"uptime_val"`        // freebox uptime (in seconds)
	BoardName        string     `json:"board_name"`        // freebox hardware revision
	TempCPUM         int        `json:"temp_cpum"`         // temp cpum (°C)
	TempSW           int        `json:"temp_sw"`           // temp sw (°C)
	TempCPUB         int        `json:"temp_cpub"`         // temp cpub (°C)
	FanRPM           int        `json:"fan_rpm"`           // fan rpm
	BoxAuthenticated bool       `json:"box_authenticated"` // is the box authenticated
	DiskStatus       diskStatus `json:"disk_status"`       // internal disk status
	BoxFlavor        boxFlavor  `json:"box_flavor"`        // box flavor
	UserMainStorage  string     `json:"user_main_storage"` // label of the storage partition for user data
}
