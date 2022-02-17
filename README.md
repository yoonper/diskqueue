# DiskQueue

Disk-based FIFO queue

---

## Features

- FIFO
- High performance

---

## Getting Started

```
go get -u github.com/yoonper/diskqueue
```

```
package main

import (
	"diskqueue"
	"fmt"
)

func main() {
	// Init
	diskqueue.Config.Path = "/tmp"
	queue := diskqueue.Init()

	// write data
	err := queue.Write([]byte("data"))
	fmt.Println(err)

	// read data
	if data, err := queue.Read(); err != nil {
		fmt.Println(data)
	}
}
```

---

## Default Config

```
Config = &config{
	Path:           "data",           // data path
	BatchSize:      100,              // number per sync
	BatchTime:      time.Second,      // interval per sync
	SegmentPerm:    0600,             // segment's mode and permission bits
	SegmentSize:    50 * 1024 * 1024, // size of each segment
	CheckpointFile: ".checkpoint",    // record read offset
}
```
