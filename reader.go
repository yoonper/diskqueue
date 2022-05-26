package diskqueue

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
)

type reader struct {
	file   *os.File
	offset int64
	reader *bufio.Reader
}

// read data
func (r *reader) read() ([]byte, error) {
	// open a new segment
	if err := r.open(); err != nil {
		return nil, err
	}

	// read a line
	data, _, err := r.reader.ReadLine()
	if err != nil {
		return nil, err
	}

	r.offset += int64(len(data)) + 1
	return data, err
}

// open a new segment
func (r *reader) open() error {
	if r.file != nil {
		return nil
	}

	files, err := r.list()
	if err != nil {
		return err
	}

	// open the earliest segment
	if r.file, err = os.OpenFile(files[0], os.O_RDONLY, Config.FilePerm); err != nil {
		return err
	}

	// seek read offset
	if _, err = r.file.Seek(r.offset, 0); err != nil {
		return err
	}

	r.reader = bufio.NewReader(r.file)
	return err
}

// rotate to next segment
func (r *reader) rotate() error {
	if r.file == nil {
		return nil
	}

	// close segment
	r.file.Close()

	// remove segment
	if err := os.Remove(r.file.Name()); err != nil {
		return err
	}

	r.file, r.offset, r.reader = nil, 0, nil
	r.sync()

	return nil
}

// list all segments
func (r *reader) list() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(Config.Path, "*.data"))
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.New("empty")
	}

	sort.Strings(files)
	return files, err
}

// sync read offset
func (r *reader) sync() {
	name := path.Join(Config.Path, Config.CheckpointFile)
	offset := []byte(strconv.FormatInt(r.offset, 10))
	_ = ioutil.WriteFile(name, offset, Config.FilePerm)
}

// close reader
func (r *reader) close() {
	if r.file == nil {
		return
	}

	r.sync()
	if err := r.file.Close(); err != nil {
		return
	}

	r.file, r.offset, r.reader = nil, 0, nil
}

// restore read offset
func (r *reader) restore() {
	name := path.Join(Config.Path, Config.CheckpointFile)
	offset, _ := ioutil.ReadFile(name)
	r.offset, _ = strconv.ParseInt(string(offset), 10, 64)
}
