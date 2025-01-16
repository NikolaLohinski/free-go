package types

import (
	"io"
)

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

type fileTaskType string

const (
	FileTaskTypeConcatenate fileTaskType = "cat"     // Concatenate multiple files
	FileTaskTypeCopy        fileTaskType = "cp"      // Copy files
	FileTaskTypeMove        fileTaskType = "mv"      // Move files
	FileTaskTypeRemove      fileTaskType = "rm"      // Remove files
	FileTaskTypeArchive     fileTaskType = "archive" // Creates an archive
	FileTaskTypeExtract     fileTaskType = "extract" // Extract an archive
	FileTaskTypeRepair      fileTaskType = "repair"  // Check and repair files

	// Undocumented and reverse engineered task types
	FileTaskTypeHash fileTaskType = "hash" // Hash a file
)

type fileTaskState string

const (
	FileTaskStateQueued  fileTaskState = "queued"  // Queued (only one task is active at a given time)
	FileTaskStateRunning fileTaskState = "running" // Running
	FileTaskStatePaused  fileTaskState = "paused"  // Paused (user suspended)
	FileTaskStateDone    fileTaskState = "done"    // Done
	FileTaskStateFailed  fileTaskState = "failed"  // Failed (see error)
)

type fileTaskError string

const (
	FileTaskErrorNone                fileTaskError = "none"                  //  No error
	FileTaskErrorArchiveReadFailed   fileTaskError = "archive_read_failed"   //  Error reading archive
	FileTaskErrorArchiveOpenFailed   fileTaskError = "archive_open_failed"   //  Error opening archive
	FileTaskErrorArchiveWriteFailed  fileTaskError = "archive_write_failed"  //  Error writing archive
	FileTaskErrorChdirFailed         fileTaskError = "chdir_failed"          //  Error changing directory
	FileTaskErrorDestIsNotDir        fileTaskError = "dest_is_not_dir"       //  The destination is not a directory
	FileTaskErrorFileExists          fileTaskError = "file_exists"           //  File already exists
	FileTaskErrorFileNotFound        fileTaskError = "file_not_found"        //  File not found
	FileTaskErrorMkdirFailed         fileTaskError = "mkdir_failed"          //  Unable to create directory
	FileTaskErrorOpenInputFailed     fileTaskError = "open_input_failed"     //  Error opening input file
	FileTaskErrorOpenOutputFailed    fileTaskError = "open_output_failed"    //  Error opening output file
	FileTaskErrorOpenDirFailed       fileTaskError = "opendir_failed"        //  Error opening directory
	FileTaskErrorOverwriteFailed     fileTaskError = "overwrite_failed"      //  Error overwriting file
	FileTaskErrorPathTooBig          fileTaskError = "path_too_big"          //  Path is too long
	FileTaskErrorRepairFailed        fileTaskError = "repair_failed"         //  Failed to repair corrupted files
	FileTaskErrorRmdirFailed         fileTaskError = "rmdir_failed"          //  Error removing directory
	FileTaskErrorSameFile            fileTaskError = "same_file"             //  Source and Destination are the same file
	FileTaskErrorUnlinkFailed        fileTaskError = "unlink_failed"         //  Error removing file
	FileTaskErrorUnsupportedFileType fileTaskError = "unsupported_file_type" //  This file type is not supported
	FileTaskErrorWriteFailed         fileTaskError = "write_failed"          //  Error writing file
	FileTaskErrorDiskFull            fileTaskError = "disk_full"             //  Disk is full
	FileTaskErrorInternal            fileTaskError = "internal"              //  Internal error
	FileTaskErrorInvalidFormat       fileTaskError = "invalid_format"        //  Invalid file format (corrupted ?)
	FileTaskErrorIncorrectPassword   fileTaskError = "incorrect_password"    //  Invalid or missing password for extraction
	FileTaskErrorPermissionDenied    fileTaskError = "permission_denied"     //  Permission denied
	FileTaskErrorReadlinkFailed      fileTaskError = "readlink_failed"       //  Failed to read the target of a symbolic link
	FileTaskErrorSymlinkFailed       fileTaskError = "symlink_failed"        //  Failed to create a symbolic link
	FileTaskErrorCopyIntoItself      fileTaskError = "copy_into_itself"      //  Attempted to copy a directory to a subdirectory of itself
	FileTaskErrorTruncateFailed      fileTaskError = "truncate_failed"       //  Failed to truncate file
	FileTaskErrorUnknownHashType     fileTaskError = "unknown_hash_type"     //  Invalid hash type
	FileTaskErrorInvalidID           fileTaskError = "invalid_id"            //  Invalid ID
)

type FileSystemTask struct {
	ID                            int64         `json:"id"`
	Type                          fileTaskType  `json:"type"`
	State                         fileTaskState `json:"state"`
	Error                         fileTaskError `json:"error"`
	CurrentBytesDone              int64         `json:"curr_bytes_done"`
	TotalBytes                    int64         `json:"total_bytes"`
	NumberFilesDone               int64         `json:"nfiles_done"`
	StartedTimestamp              int64         `json:"started_ts"`
	DurationSeconds               int64         `json:"duration"`
	DoneTimestamp                 int64         `json:"done_ts"`
	CurrentBytes                  int64         `json:"curr_bytes"`
	To                            string        `json:"to"`
	NumberFiles                   int64         `json:"nfiles"`
	CreatedTimestamp              int64         `json:"created_ts"`
	TotalBytesDone                int64         `json:"total_bytes_done"`
	From                          string        `json:"from"`
	ProcessingRate                int64         `json:"rate"`
	EstimatedTimeRemainingSeconds int64         `json:"eta"`
	ProgressPercent               int           `json:"progress"`
	Sources                       []string      `json:"src"`
	Destination                   string        `json:"dst"`
}

type HashType string

const (
	HashTypeMD5    HashType = "md5"
	HashTypeSHA1   HashType = "sha1"
	HashTypeSHA256 HashType = "sha256"
	HashTypeSHA512 HashType = "sha512"
)

type HashPayload struct {
	HashType HashType   `json:"hash_type"`
	Path     Base64Path `json:"src"`
}

type HashResult struct {
	Hash string `json:"hash"`
}

type File struct {
	ContentType string
	FileName 	string
	Content     io.Reader
}

type FileSytemTaskUpdate struct {
	State fileTaskState `json:"state"`
}

type FileMoveMode string

const (
	FileMoveModeOverwrite FileMoveMode = "overwrite" // Overwrite the destination file
	FileMoveModeSkip      FileMoveMode = "skip"      // Keep the destination file
	FileMoveModeBoth      FileMoveMode = "both"      // Keep both files (rename the file adding a suffix)
	FileMoveModeRecent    FileMoveMode = "recent"    // Only overwrite if newer than destination file
)

type FileCopyMode string

const (
	FileCopyModeOverwrite FileCopyMode = "overwrite" // Overwrite the destination file
	FileCopyModeSkip      FileCopyMode = "skip"      // Keep the destination file
	FileCopyModeBoth      FileCopyMode = "both"      // Keep both files (rename the file adding a suffix)
	FileCopyModeRecent    FileCopyMode = "recent"    // Only overwrite if newer than destination file
)
