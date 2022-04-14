package diskqueue

import (
	"errors"
	"io"
	"sync"
	"time"
)

type Diskqueue struct {
	sync.RWMutex
	close  bool
	ticker *time.Ticker
}

var (
	Writer = &writer{}
	Reader = &reader{}

	Config = &config{
		Path:              "data",
		FilePerm:          0600,
		BatchSize:         100,
		BatchTime:         time.Second,
		SegmentSize:       50 * 1024 * 1024,
		SegmentLimit:      2048,
		CheckpointFile:    ".checkpoint",
		MiniRequiredSpace: 1024 * 1024 * 1024,
	}

	ErrEmpty  = errors.New("empty")
	ErrClosed = errors.New("closed")
)

// Init diskqueue
func Init() *Diskqueue {
	queue := &Diskqueue{close: false}
	queue.ticker = time.NewTicker(Config.BatchTime)
	Reader.restore()

	go func() {
		for {
			<-queue.ticker.C
			queue.Lock()
			Writer.sync()
			Reader.sync()
			queue.Unlock()
		}
	}()
	return queue
}

// Write data
func (queue *Diskqueue) Write(data []byte) error {
	if queue.close {
		return ErrClosed
	}

	queue.Lock()
	defer queue.Unlock()

	return Writer.write(data)
}

// Read data
func (queue *Diskqueue) Read() ([]byte, error) {
	if queue.close {
		return nil, ErrClosed
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
	queue.close = true
	Writer.sync()
	Reader.sync()
}
