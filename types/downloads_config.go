package types

type DlThrottlingMode string

const (
	DlThrottlingModeNormal    DlThrottlingMode = "normal"
	DlThrottlingModeSlow      DlThrottlingMode = "slow"
	DlThrottlingModeHibernate DlThrottlingMode = "hibernate"
	DlThrottlingModeSchedule  DlThrottlingMode = "schedule"
)

type DlRate struct {
	TxRate int64 `json:"tx_rate"` // maximum transmit rate (in byte/s), 0 means no limit
	RxRate int64 `json:"rx_rate"` // maximum receive rate (in byte/s), 0 means no limit
}

type DlThrottlingConfig struct {
	Normal   DlRate             `json:"normal"`
	Slow     DlRate             `json:"slow"`
	Schedule []DlThrottlingMode `json:"schedule"` // 168-element array (24*7 week hours starting Monday midnight)
	Mode     DlThrottlingMode   `json:"mode"`
}

type DlBtCryptoSupport string

const (
	DlBtCryptoSupportUnsupported DlBtCryptoSupport = "unsupported"
	DlBtCryptoSupportAllowed     DlBtCryptoSupport = "allowed"
	DlBtCryptoSupportPreferred   DlBtCryptoSupport = "preferred"
	DlBtCryptoSupportRequired    DlBtCryptoSupport = "required"
)

type DlBtConfig struct {
	MaxPeers        int               `json:"max_peers"`
	StopRatio       int               `json:"stop_ratio"` // scaled by 100; 150 means 1.5x
	CryptoSupport   DlBtCryptoSupport `json:"crypto_support"`
	EnableDHT       bool              `json:"enable_dht"`
	EnablePex       bool              `json:"enable_pex"`
	AnnounceTimeout int               `json:"announce_timeout"` // seconds
	MainPort        int               `json:"main_port"`
	DHTPort         int               `json:"dht_port"`
}

type DlNewsConfig struct {
	Server      string `json:"server"`
	Port        int    `json:"port"`
	SSL         bool   `json:"ssl"`
	User        string `json:"user"`
	Password    string `json:"password,omitempty"` // write-only
	NThreads    int    `json:"nthreads"`
	AutoRepair  bool   `json:"auto_repair"`
	LazyPar2    bool   `json:"lazy_par2"`
	AutoExtract bool   `json:"auto_extract"`
	EraseTmp    bool   `json:"erase_tmp"`
}

type DlFeedConfig struct {
	FetchInterval int `json:"fetch_interval"` // minutes
	MaxItems      int `json:"max_items"`
}

type DlBlockListConfig struct {
	Sources interface{} `json:"sources"`
}

type DownloadConfiguration struct {
	MaxDownloadingTasks int                `json:"max_downloading_tasks"`
	DownloadDir         Base64Path         `json:"download_dir"`
	WatchDir            Base64Path         `json:"watch_dir"`
	UseWatchDir         bool               `json:"use_watch_dir"`
	Throttling          DlThrottlingConfig `json:"throttling"`
	News                DlNewsConfig       `json:"news"`
	Bt                  DlBtConfig         `json:"bt"`
	Feed                DlFeedConfig       `json:"feed"`
	BlockList           DlBlockListConfig  `json:"blocklist"`
	DNS1                string             `json:"dns1"`
	DNS2                string             `json:"dns2"`
}
