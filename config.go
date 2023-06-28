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
	WriteTimeout     int64         // close segment if timeout (in seconds)
	CheckpointFile   string        // read index and offset
	MinRequiredSpace int64         // minimum required free space (in bytes)
}
