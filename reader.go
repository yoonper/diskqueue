package diskqueue

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
)

type reader struct {
	file       *os.File
	index      int64
	offset     int64
	reader     *bufio.Reader
	checkpoint checkpoint
}

type checkpoint struct {
	Index  int64 `json:"index"`
	Offset int64 `json:"offset"`
}

// read data
func (r *reader) read() (int64, int64, []byte, error) {
	if err := r.check(); err != nil {
		return r.index, r.offset, nil, err
	}

	// read a line
	data, _, err := r.reader.ReadLine()
	if err != nil {
		return r.index, r.offset, nil, err
	}

	r.offset += int64(len(data)) + 1
	return r.index, r.offset, data, err
}

// check a new segment
func (r *reader) check() error {
	if r.file != nil {
		return nil
	}

	file, err := r.next()
	if err != nil {
		return err
	}

	return r.open(file)
}

func (r *reader) open(file string) (err error) {
	if r.file, err = os.OpenFile(file, os.O_RDONLY, Config.FilePerm); err != nil {
		return err
	}

	// get file index
	r.index = r.getIndex(file)

	// seek read offset
	if _, err = r.file.Seek(r.offset, 0); err != nil {
		return err
	}

	r.reader = bufio.NewReader(r.file)
	return nil
}

// rotate to next segment
func (r *reader) rotate() error {
	if r.file == nil {
		return nil
	}

	// close segment
	_ = r.file.Close()
	r.file, r.offset, r.reader = nil, 0, nil
	return nil
}

// close reader
func (r *reader) close() {
	if r.file == nil {
		return
	}

	if err := r.file.Close(); err != nil {
		return
	}

	r.file, r.reader, r.index, r.offset = nil, nil, 0, 0
}

// sync index and offset
func (r *reader) sync() {
	name := path.Join(Config.Path, Config.CheckpointFile)
	data, _ := json.Marshal(&r.checkpoint)
	_ = ioutil.WriteFile(name, data, Config.FilePerm)
}

// restore index and offset
func (r *reader) restore() (err error) {
	name := path.Join(Config.Path, Config.CheckpointFile)
	data, _ := ioutil.ReadFile(name)
	_ = json.Unmarshal(data, &r.checkpoint)
	r.index, r.offset = r.checkpoint.Index, r.checkpoint.Offset

	if r.index == 0 {
		return nil
	}

	if err = r.open(fmt.Sprintf("%s/%d.data", Config.Path, r.index)); err != nil {
		r.offset = 0
	}

	return err
}

// next segment
func (r *reader) next() (string, error) {
	files, err := filepath.Glob(filepath.Join(Config.Path, "*.data"))
	if err != nil {
		return "", err
	}

	sort.Strings(files)

	for _, file := range files {
		index := r.getIndex(file)
		if index < r.checkpoint.Index {
			_ = os.Remove(file) // remove expired segment
		}

		if index > r.index {
			return file, nil
		}
	}

	return "", errors.New("queue is empty")
}

// get segment index
func (r *reader) getIndex(filename string) int64 {
	base := path.Base(filename)
	name := base[0 : len(base)-len(path.Ext(filename))]
	index, _ := strconv.ParseInt(name, 10, 64)
	return index
}
