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
	count  int64
	writer *bufio.Writer
}

// write data
func (w *writer) write(data []byte) error {
	// append newline
	data = append(data, "\n"...)
	size := int64(len(data))

	// close current segment for rotate
	if w.size+size > Config.SegmentSize {
		w.close()
	}

	// create a new segment
	if w.file == nil {
		if err := w.open(); err != nil {
			return err
		}
	}

	// write to buffer
	if _, err := w.writer.Write(data); err != nil {
		return err
	}

	w.size += size

	// sync data to disk
	w.count++
	if w.count >= Config.BatchSize {
		w.sync()
	}

	return nil
}

// create a new segment
func (w *writer) open() error {
	var err error
	name := path.Join(Config.Path, fmt.Sprintf("%013d.data", time.Now().UnixNano()/1e6))
	if w.file, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY, Config.FilePerm); err != nil {
		return err
	}

	w.size = 0
	// disable auto flush
	w.writer = bufio.NewWriterSize(w.file, int(Config.SegmentSize))
	return err
}

// sync data to disk
func (w *writer) sync() {
	if w.writer == nil {
		return
	}

	if err := w.writer.Flush(); err == nil {
		w.count = 0
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
