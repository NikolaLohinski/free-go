package types

type APIVersion struct {
	UID            string `json:"uid"`
	DeviceName     string `json:"device_name"`
	DeviceType     string `json:"device_type"`
	APIVersion     string `json:"api_version"`
	APIDomain      string `json:"api_domain"`
	APIBaseURL     string `json:"api_base_url"`
	BoxModelName   string `json:"box_model_name"`
	BoxModel       string `json:"box_model"`
	HTTPSPort      int    `json:"https_port"`
	HTTPSAvailable bool   `json:"https_available"`
}
