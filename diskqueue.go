package diskqueue

import (
	"errors"
	"io"
	"os"
	"sync"
	"time"
)

type Diskqueue struct {
	sync.RWMutex
	close     bool
	closeChan chan bool
	ticker    *time.Ticker
}

var (
	Writer = &writer{}
	Reader = &reader{}

	Config = &config{
		Path:             "data",
		FilePerm:         0600,
		BatchSize:        100,
		BatchTime:        time.Second,
		SegmentSize:      50 * 1024 * 1024,
		SegmentLimit:     2048,
		CheckpointFile:   ".checkpoint",
		MinRequiredSpace: 1024 * 1024 * 1024,
	}
)

// Start diskqueue
func Start() (*Diskqueue, error) {
	if _, err := os.Stat(Config.Path); err != nil {
		return nil, err
	}

	queue := &Diskqueue{close: false, closeChan: make(chan bool)}
	queue.ticker = time.NewTicker(Config.BatchTime)
	Reader.restore()
	go queue.sync()
	return queue, nil
}

// Write data
func (queue *Diskqueue) Write(data []byte) error {
	if queue.close {
		return errors.New("closed")
	}

	queue.Lock()
	defer queue.Unlock()

	return Writer.write(data)
}

// Read data
func (queue *Diskqueue) Read() ([]byte, error) {
	if queue.close {
		return nil, errors.New("closed")
	}

	queue.RLock()
	defer queue.RUnlock()

	data, err := Reader.read()
	if err == io.EOF && (Writer.file == nil || Reader.file.Name() != Writer.file.Name()) {
		_ = Reader.rotate()
	}
	return data, err
}

// Close diskqueue
func (queue *Diskqueue) Close() {
	if queue.close {
		return
	}

	queue.close = true
	queue.closeChan <- true
	Writer.close()
	Reader.close()
}

// sync data
func (queue *Diskqueue) sync() {
	for {
		select {
		case <-queue.ticker.C:
			queue.Lock()
			Writer.sync()
			Reader.sync()
			queue.Unlock()
		case <-queue.closeChan:
			return
		}
	}
}
