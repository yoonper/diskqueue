package diskqueue

import (
	"os"
	"time"
)

type config struct {
	Path           string        // data path
	FilePerm       os.FileMode   // file's mode and permission bits
	BatchSize      int64         // number per sync
	BatchTime      time.Duration // interval per sync
	SegmentSize    int64         // size of each segment
	CheckpointFile string        // record read offset
}
