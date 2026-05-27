package types

type SambaConfiguration struct {
	FileShareEnabled  bool   `json:"file_share_enabled"`  // is file sharing enabled
	PrintShareEnabled bool   `json:"print_share_enabled"` // is printer sharing enabled
	LogonEnabled      bool   `json:"logon_enabled"`       // is login/password required to access shares
	LogonUser         string `json:"logon_user"`          // login user name
	Workgroup         string `json:"workgroup"`           // workgroup name
	V2Enabled         bool   `json:"smbv2_enabled"`       // Set to true to enable SMBv2/v3
}

type SambaConfigurationPayload struct {
	SambaConfiguration `json:",inline"`
	LoginPassword      string `json:"login_password"` // Samba user password
}

type AFPConfiguration struct {
	Enabled    bool                  `json:"enabled"`     // is afp service enabled
	GuestAllow string                `json:"guest_allow"` // allow guest to access shared files
	ServerType netshareAFPServerType `json:"server_type"` // Afp server type (to display proper icon) in MacOS
	LoginName  string                `json:"login_name"`  // Afp user name
}

type AFPConfigurationPayload struct {
	AFPConfiguration `json:",inline"`
	LoginPassword    string `json:"login_password"` // Afp user password
}

type netshareAFPServerType = string

const (
	NetshareAFPServerTypePowerBook  netshareAFPServerType = "powerbook"
	NetshareAFPServerTypePowerMac   netshareAFPServerType = "powermac"
	NetshareAFPServerTypeMacMini    netshareAFPServerType = "macmini"
	NetshareAFPServerTypeIMac       netshareAFPServerType = "imac"
	NetshareAFPServerTypeMacBook    netshareAFPServerType = "macbook"
	NetshareAFPServerTypeMacBookPro netshareAFPServerType = "macbookpro"
	NetshareAFPServerTypeMacBookAir netshareAFPServerType = "macbookair"
	NetshareAFPServerTypeMacPro     netshareAFPServerType = "macpro"
	NetshareAFPServerTypeAppleTV    netshareAFPServerType = "appletv"
	NetshareAFPServerTypeAirport    netshareAFPServerType = "airport"
	NetshareAFPServerTypeXServe     netshareAFPServerType = "xserve"
)

type netshareError string

const (
	NetshareErrorInvalidWorkgroupName    netshareError = "invalid_workgroup_name"     // Invalid workgroup name
	NetshareErrorInvalidLogonUser        netshareError = "invalid_logon_user"         // Invalid samba user name
	NetshareErrorInvalidLogonPassword    netshareError = "invalid_logon_password"     // Invalid samba user password
	NetshareErrorInvalidAFPLoginName     netshareError = "invalid_afp_login_name"     // Invalid AFP user name
	NetshareErrorInvalidAFPLoginPassword netshareError = "invalid_afp_login_password" // Invalid AFP user password
)
