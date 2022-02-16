package diskqueue

import (
	"os"
	"time"
)

type config struct {
	Path           string        // data path
	BatchTime      time.Duration // interval per sync
	SegmentPerm    os.FileMode   // segment's mode and permission bits
	SegmentSize    int64         // size of each segment
	CheckpointFile string        // record read offset
}
