package types

type fileType string

const (
	FileTypeDirectory fileType = "dir"
	FileTypeFile      fileType = "file"
)

type FileInfo struct {
	Type         fileType   `json:"type"`
	Index        int64      `json:"index"`
	Link         bool       `json:"link"`
	Parent       Base64Path `json:"parent"`
	Modification uint64     `json:"modification"`
	Hidden       bool       `json:"hidden"`
	MimeType     string     `json:"mimetype"`
	Name         string     `json:"name"`
	Path         Base64Path `json:"path"`
	SizeBytes    uint64     `json:"size"`
}
