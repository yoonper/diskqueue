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
	// init
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
	FilePerm:       0600,             // file's mode and permission bits
	BatchSize:      100,              // number per sync
	BatchTime:      time.Second,      // interval per sync
	SegmentSize:    50 * 1024 * 1024, // size of each segment
	CheckpointFile: ".checkpoint",    // record read offset
}
```
