package kodik

// Error codes for kodik CLI, matching Bash script semantics.
const (
	ErrNetwork    = 10 // Network error
	ErrChecksum   = 20 // Checksum validation error
	ErrFilesystem = 30 // Filesystem error
	ErrUser       = 40 // User error
	ErrSystem     = 50 // System error
)

// KodikError wraps an error with a code and message.
type KodikError struct {
	Code    int
	Message string
}

func (e *KodikError) Error() string {
	return e.Message
}

func NewKodikError(code int, msg string) *KodikError {
	return &KodikError{Code: code, Message: msg}
}

// PreservationError provides rich context when a preservation-aware operation fails.
// It distinguishes what was preserved vs. what operation failed, and where backups live.
type PreservationError struct {
	Operation  string   // e.g. "mergeGithubComponent"
	Component  string   // e.g. ".github"
	UserFiles  []string // user-managed files retained
	KodikFiles []string // kodik-managed files involved in the operation
	BackupPath string   // path to backup directory
	Cause      error    // underlying error
}

func (e *PreservationError) Error() string {
	msg := "preservation failure"
	if e.Operation != "" {
		msg += " op=" + e.Operation
	}
	if e.Component != "" {
		msg += " component=" + e.Component
	}
	if e.BackupPath != "" {
		msg += " backup=" + e.BackupPath
	}
	if e.Cause != nil {
		msg += ": " + e.Cause.Error()
	}
	return msg
}

// NewPreservationError helper for constructing a PreservationError.
func NewPreservationError(op, comp, backup string, userFiles, kodikFiles []string, cause error) *PreservationError {
	return &PreservationError{
		Operation:  op,
		Component:  comp,
		UserFiles:  userFiles,
		KodikFiles: kodikFiles,
		BackupPath: backup,
		Cause:      cause,
	}
}
