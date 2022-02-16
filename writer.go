package diskqueue

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"time"
)

type writer struct {
	file   *os.File
	size   int64
	writer *bufio.Writer
}

// write data
func (w *writer) write(data []byte) error {
	// create a new segment
	if w.file == nil {
		if err := w.open(); err != nil {
			return err
		}
	}

	// append a newline and write to buffer
	if _, err := w.writer.Write(append(data, "\n"...)); err != nil {
		return err
	}

	// close current segment for rotate
	w.size += int64(len(data))
	if w.size > Config.SegmentSize {
		w.close()
	}

	return nil
}

// create a new segment
func (w *writer) open() error {
	var err error
	name := path.Join(Config.Path, fmt.Sprintf("%013d.data", time.Now().UnixNano()/1e6))
	if w.file, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY, Config.SegmentPerm); err != nil {
		return err
	}

	w.size = 0
	//w.writer = bufio.NewWriter(w.file)
	w.writer = bufio.NewWriterSize(w.file, 99999999999999999)
	return err
}

// sync data to disk
func (w *writer) sync() {
	if w.writer != nil {
		_ = w.writer.Flush()
	}
}

// close segment
func (w *writer) close() {
	if w.file == nil {
		return
	}

	w.sync()
	if err := w.file.Close(); err != nil {
		return
	}

	w.size, w.file, w.writer = 0, nil, nil
}
