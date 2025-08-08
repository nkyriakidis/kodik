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
