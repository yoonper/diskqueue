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
	"fmt"
	"github.com/yoonper/diskqueue"
	"log"
	"time"
)

func main() {
	var err error
	var queue *diskqueue.Diskqueue

	// config
	diskqueue.Config.Path = "/tmp/diskqueue"
	diskqueue.Config.BatchSize = 1

	// start
	if queue, err = diskqueue.Start(); err != nil {
		log.Fatalln(err)
	}

	// write
	go func() {
		for {
			time.Sleep(time.Second)
			data := []byte(time.Now().Format("2006-01-02 15:04:05"))
			if err := queue.Write(data); err != nil {
				fmt.Println(err)
			}
		}
	}()

	// read
	go func() {
		for {
			time.Sleep(time.Second)
			if index, offset, data, err := queue.Read(); err == nil {
				fmt.Println(index, offset, string(data))
				queue.Commit(index, offset) // commit
			}
		}
	}()

	select {}
}
```

---

## Default Config

```
Config = &config{
	Path:              "data",
	FilePerm:          0600,
	BatchSize:         100,
	BatchTime:         time.Second,
	SegmentSize:       50 * 1024 * 1024,
	SegmentLimit:      2048,
	WriteTimeout       300,
	CheckpointFile:    ".checkpoint",
	MinRequiredSpace:  1024 * 1024 * 1024,
}
```
