package core

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// LogWriter writes to a log file and optionally calls a line callback.
type LogWriter struct {
	file     *os.File
	mu       sync.Mutex
	callback func(string)
}

// NewLogWriter creates a LogWriter writing to the given path.
func NewLogWriter(logDir string) (*LogWriter, error) {
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(logDir, "core.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &LogWriter{file: f}, nil
}

// SetCallback sets the line callback for real-time streaming.
func (w *LogWriter) SetCallback(cb func(string)) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.callback = cb
}

// Write implements io.Writer. It writes to the log file and calls the
// line callback for each line of output.
func (w *LogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	cb := w.callback
	w.mu.Unlock()

	if cb != nil {
		scanner := bufio.NewScanner(newBytesReader(p))
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				cb(line)
			}
		}
	}

	if w.file != nil {
		return w.file.Write(p)
	}
	return len(p), nil
}

// Close closes the underlying file.
func (w *LogWriter) Close() error {
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// newBytesReader is a helper to create an io.Reader from a byte slice.
func newBytesReader(p []byte) io.Reader {
	return &bytesReader{data: p}
}

type bytesReader struct {
	data []byte
	pos  int
}

func (r *bytesReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
