package log

import (
	"fmt"
	"os"
	"sync"
)

const logRotateNumberOfFiles = 3

// maximum for this value is max(int32)
const logRotateFileByteMaximum = 2e7

// RotateWriter represents io.Writer that writes to file with rotation
type RotateWriter struct {
	Lock     sync.Mutex
	filename string
	Fp       *os.File
	fileSize int
}

// NewRotateWriter initializes RotateWriter
func NewRotateWriter(filename string) (*RotateWriter, error) {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	return &RotateWriter{
		filename: filename,
		fileSize: int(fileInfo.Size()),
		Fp:       fp,
	}, nil
}

func (w *RotateWriter) Write(output []byte) (int, error) {
	w.Lock.Lock()
	defer w.Lock.Unlock()

	bytesWritten, err := w.Fp.Write(output)
	if err != nil {
		return bytesWritten, err
	}
	w.fileSize += bytesWritten
	if w.fileSize > logRotateFileByteMaximum {
		err = w.rotate()
		if err != nil {
			return bytesWritten, err
		}
	}
	return bytesWritten, nil
}

func (w *RotateWriter) rotate() error {
	// this method is called when mutex is already locked

	// Close opened file descriptor stored in global variable
	err := w.Fp.Close()
	if err != nil {
		return err
	}

	// perform file rotation
	for i := logRotateNumberOfFiles - 2; i >= 0; i-- {
		var sourceName string
		if i != 0 {
			sourceName = fmt.Sprintf("%s.%d", w.filename, i)
		} else {
			sourceName = w.filename
		}
		destinationName := fmt.Sprintf("%s.%d", w.filename, i+1)
		// move sourceName to DestinationName if sourceName exists
		if _, err = os.Stat(sourceName); err == nil {
			err = os.Rename(sourceName, destinationName)
			if err != nil {
				return err
			}
		}
	}

	// Open file and store descriptor in global variable
	fp, err := os.OpenFile(w.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	w.Fp = fp

	// Reset bytes counter
	w.fileSize = 0
	return err
}
