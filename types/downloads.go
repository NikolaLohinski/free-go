package types

type downloadTaskType string

const (
	DownloadTaskTypeBitTorrent downloadTaskType = "bt"   // bittorrent download
	DownloadTaskTypeNewsGroup  downloadTaskType = "nzb"  // newsgroup download
	DownloadTaskTypeHTTP       downloadTaskType = "http" // HTTP download
	DownloadTaskTypeFTP        downloadTaskType = "ftp"  // FTP download
)

type downloadTaskStatus string

const (
	DownloadTaskStatusStopped     = "stopped"    //	task is stopped, can be resumed by setting the status to downloading
	DownloadTaskStatusQueued      = "queued"     //	task will start when a new download slot is available the queue position is stored in queue_pos attribute
	DownloadTaskStatusStarting    = "starting"   //	task is preparing to start download
	DownloadTaskStatusStopping    = "stopping"   //	task is gracefully stopping
	DownloadTaskStatusError       = "error"      //	there was a problem with the download, you can get an error code in the error field
	DownloadTaskStatusDone        = "done"       //	the download is over. For bt you can resume seeding setting the status to seeding if the ratio is not reached yet
	DownloadTaskStatusChecking    = "checking"   //	(only valid for nzb) download is over, the downloaded files are being checked using par2
	DownloadTaskStatusRepairing   = "repairing"  //	(only valid for nzb) download is over, the downloaded files are being repaired using par2
	DownloadTaskStatusExtracting  = "extracting" //	only valid for nzb) download is over, the downloaded files are being extracted
	DownloadTaskStatusSeeding     = "seeding"    //	(only valid for bt) download is over, the content is Change to being shared to other users. The task will automatically stop once the seed ratio has been reached
	DownloadTaskStatusRetry       = "retry"      //	You can set a task status to ‘retry’ to restart the download task.
	DownloadTaskStatusDownloading = "downloading"
)

type downloadTaskIOPriority string

const (
	DownloadTaskIOPriorityLow    downloadTaskIOPriority = "low"
	DownloadTaskIOPriorityNormal downloadTaskIOPriority = "normal"
	DownloadTaskIOPriorityHigh   downloadTaskIOPriority = "high"
)

type DownloadTask struct {
	ID                 int64                  `json:"id"`
	Type               downloadTaskType       `json:"type"`
	Name               string                 `json:"name"`
	Status             downloadTaskStatus     `json:"status"`
	IOPriority         downloadTaskIOPriority `json:"io_priority"`
	SizeBytes          int64                  `json:"size"`             // Download size (in Bytes)
	QueuePosition      int64                  `json:"queue_pos"`        // position in download queue (0 if not queued)
	TransmittedBytes   int64                  `json:"tx_bytes"`         // transmitted bytes (including protocol overhead)
	ReceivedBytes      int64                  `json:"rx_bytes"`         // received bytes (including protocol overhead)
	TransmitRate       int64                  `json:"tx_rate"`          // current transmit rate (in byte/s)
	ReceiveRate        int64                  `json:"rx_rate"`          // current receive rate (in byte/s)
	TransmitPercentage int                    `json:"tx_pct"`           // transmit percentage (without protocol overhead). To improve precision the value as been scaled by 100 so that a tx_pct of 123 means 1.23%
	ReceivedPercentage int                    `json:"rx_pct"`           // received percentage (without protocol overhead). To improve precision the value as been scaled by 100 so that a tx_pct of 123 means 1.23%
	Error              string                 `json:"error"`            // an error code
	CreatedTimestamp   Timestamp              `json:"created_ts"`       // timestamp of the download creation time
	ETASeconds         int64                  `json:"eta"`              // estimated remaining download time (in seconds)
	DownloadDirectory  Base64Path             `json:"download_dir"`     // directory where the file(s) will be saved (base64 encoded)
	StopRatio          int64                  `json:"stop_ratio"`       // Only relevant for bittorrent tasks. Once the transmit ration has been reached the task will stop seeding. The ratio is scaled by 100 to improve resolution. A stop_ratio of 150 means that the task will stop seeding once tx_bytes = 1.5 * rx_bytes.
	ArchivePassword    string                 `json:"archive_password"` // (only relevant for nzb) password for extracting downloaded archives
	InfoHash           string                 `json:"info_hash"`        // (only relevant for bt) torrent info_hash encoded in hexa
	PieceLength        int64                  `json:"piece_length"`     // (only relevant for bt) torrent piece length in bytes
}

type DownloadRequest struct {
	DownloadURLs      []string          // The URL(s)
	DownloadDirectory string            // The download destination directory (optional: will use the configuration download_dir by default)
	Filename          string            // Override the name of the destination file. Only valid with one, non-recursive download_url.
	Hash              string            // Verify the hash of the downloaded file. The format is sha256:xxxxxx or sha512:xxxxxx; or the URL of a SHA256SUMS, SHA512SUMS, -CHECKSUM or .sha256 file. Only valid with one, non-recursive download_url.
	Recursive         bool              // If true the download will be recursive
	Username          string            // Auth username (optional)
	Password          string            // Auth password (optional)
	ArchivePassword   string            // The password required to extract downloaded content (only relevant for nzb)
	Cookies           map[string]string // The http cookies (to be able to pass session cookies along with url). This is the content of the HTTP Cookie header, for example: cookie1=value1; cookie2=value2
}
