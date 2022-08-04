package diskqueue

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"
	"time"
)

type Diskqueue struct {
	sync.RWMutex
	close  bool
	ticker *time.Ticker
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
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
	queue := &Diskqueue{close: false, wg: &sync.WaitGroup{}}
	queue.ticker = time.NewTicker(Config.BatchTime)
	queue.ctx, queue.cancel = context.WithCancel(context.TODO())
	_ = Reader.restore()
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
func (queue *Diskqueue) Read() (int64, int64, []byte, error) {
	if queue.close {
		return 0, 0, nil, errors.New("closed")
	}

	queue.RLock()
	defer queue.RUnlock()

	index, offset, data, err := Reader.read()
	if err == io.EOF && (Writer.file == nil || Reader.file.Name() != Writer.file.Name()) {
		_ = Reader.rotate()
	}
	return index, offset, data, err
}

// Commit index and offset
func (queue *Diskqueue) Commit(index int64, offset int64) {
	if queue.close {
		return
	}

	ck := &Reader.checkpoint
	ck.Index, ck.Offset = index, offset
	Reader.sync()
}

// Close diskqueue
func (queue *Diskqueue) Close() {
	if queue.close {
		return
	}

	queue.close = true
	queue.cancel()
	queue.wg.Wait()
	Writer.close()
	Reader.close()
}

// sync data
func (queue *Diskqueue) sync() {
	queue.wg.Add(1)
	defer queue.wg.Done()
	for {
		select {
		case <-queue.ticker.C:
			queue.Lock()
			Writer.sync()
			queue.Unlock()
		case <-queue.ctx.Done():
			return
		}
	}
}
