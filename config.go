package diskqueue

import (
	"os"
	"time"
)

type config struct {
	Path             string        // data path
	FilePerm         os.FileMode   // segment's mode and permission bits
	BatchSize        int64         // number per sync
	BatchTime        time.Duration // interval per sync
	SegmentSize      int64         // size of each segment (in bytes)
	SegmentLimit     int64         // max number of segment
	CheckpointFile   string        // record read offset
	MinRequiredSpace int64         // minimum required free space (in bytes)
}
