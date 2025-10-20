package types

type uploadActionForce string

const (
	FileUploadStartActionForceOverwrite uploadActionForce = "overwrite"
	FileUploadStartActionForceMissing   uploadActionForce = "missing"
	FileUploadStartActionForceResume    uploadActionForce = "resume"
)

const (
	FileUploadStartActionNameUploadStart    WebSocketAction = "upload_start"
	FileUploadStartActionNameUploadFinalize WebSocketAction = "upload_finalize"
	FileUploadStartActionNameUploadData     WebSocketAction = "upload_data"
	FileUploadStartActionNameUploadCancel   WebSocketAction = "upload_cancel"
)

type UploadRequestID = WebSocketRequestID

type FileUploadStartAction struct {
	RequestID UploadRequestID   `json:"request_id,omitempty"` // optional request_id
	Action    WebSocketAction   `json:"action"`               // must be ‘upload_start’
	Size      int               `json:"size"`                 // optional file size
	Dirname   Base64Path        `json:"dirname"`              // the destination directory (encoded value)
	Filename  string            `json:"filename"`             // the destination filename
	Force     uploadActionForce `json:"force"`                // select the way conflicts are handled
}

type FileUploadStartResponse struct {
	Success   bool            `json:"success"`
	Action    WebSocketAction `json:"action"`
	RequestID UploadRequestID `json:"request_id"`
	Message   string          `json:"msg,omitempty"`
	FileSize  int             `json:"file_size,omitempty"`
	ErrorCode string          `json:"error_code,omitempty"`
}

func (r *FileUploadStartResponse) GetError() error {
	if !r.Success {
		return &WebSocketResponseError{
			ErrorCode: r.ErrorCode,
			Message:   r.Message,
		}
	}

	return nil
}

type FileUploadStartActionInput struct {
	Size     int               // optional file size
	Dirname  Base64Path        // the destination directory (encoded value)
	Filename string            // the destination filename
	Force    uploadActionForce // select the way conflicts are handled
}

type FileUploadChunkResponse struct {
	TotalLen  int  `json:"total_len"` // Target file current length
	Complete  bool `json:"complete"`  // Will be true in a reply to FileUploadFinalizeAction or FileUploadCancelAction
	Cancelled bool `json:"cancelled"` // Will be true in a reply FileUploadCancelAction
}

type FileUploadFinalize struct {
	Action    WebSocketAction `json:"action"`               // must be ‘upload_finalize’
	RequestID UploadRequestID `json:"request_id,omitempty"` // optional request_id
}

type FileUploadCancelAction struct {
	Action    WebSocketAction `json:"action"`               // must be ‘upload_cancel’
	RequestID UploadRequestID `json:"request_id,omitempty"` // optional request_id

}

type FileUploadCancelResponse struct {
	Success   bool            `json:"success"`
	Action    WebSocketAction `json:"action"`
	RequestID UploadRequestID `json:"request_id"`
	Message   string          `json:"msg,omitempty"`
	ErrorCode string          `json:"error_code,omitempty"`
}

func (r *FileUploadCancelResponse) GetError() error {
	if !r.Success {
		return &WebSocketResponseError{
			ErrorCode: r.ErrorCode,
			Message:   r.Message,
		}
	}

	return nil
}

type FileUploadFinalizeResponse struct {
	TotalLen  int  `json:"total_len"`           // Target file current length
	Complete  bool `json:"complete"`            // Will be true in a reply to FileUploadFinalizeAction or FileUploadCancelAction
	Cancelled bool `json:"cancelled,omitempty"` // Will be true in a reply FileUploadCancelAction
}

type UploadTaskStatus string

const (
	UploadTaskStatusAuthorized UploadTaskStatus = "authorized"  // Upload authorization is valid, upload has not started yet
	UploadTaskStatusInProgress UploadTaskStatus = "in_progress" // Upload in progress
	UploadTaskStatusDone       UploadTaskStatus = "done"        // Upload done
	UploadTaskStatusFailed     UploadTaskStatus = "failed"      // Upload failed
	UploadTaskStatusConflict   UploadTaskStatus = "conflict"    // Destination file conflict
	UploadTaskStatusTimeout    UploadTaskStatus = "timeout"     // Upload authorization is no longer valid
	UploadTaskStatusCancelled  UploadTaskStatus = "cancelled"   // Upload cancelled by user
)

type UploadTask struct {
	ID         int64            `json:"id"`          // Upload id
	Size       int64            `json:"size"`        // Upload file size in bytes
	Uploaded   int64            `json:"uploaded"`    // Uploaded bytes
	Status     UploadTaskStatus `json:"status"`      // Upload status
	StartDate  Timestamp        `json:"start_date"`  // Upload start date
	LastUpdate Timestamp        `json:"last_update"` // Last update of file upload object
	UploadName string           `json:"upload_name"` // Name of the file uploaded
	Dirname    string           `json:"dirname"`     // Upload destination directory
}
