package types

type Permissions struct {
	Parental   bool `json:"parental"`
	Player     bool `json:"player"`
	Explorer   bool `json:"explorer"`
	TV         bool `json:"tv"`
	Wdo        bool `json:"wdo"`
	Downloader bool `json:"downloader"`
	Profile    bool `json:"profile"`
	Camera     bool `json:"camera"`
	Settings   bool `json:"settings"`
	Calls      bool `json:"calls"`
	Home       bool `json:"home"`
	PVR        bool `json:"pvr"`
	VM         bool `json:"vm"`
	Contacts   bool `json:"contacts"`
}
