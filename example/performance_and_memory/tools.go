package main

import (
	"encoding/json"
	"fmt"
	"github.com/akley-MK4/net-defragmenter/stats"
	PCI "github.com/akley-MK4/pep-coroutine/implement"
	"log"
	"runtime"
)

type memorySnapshot struct {
	Title      string
	AllocBytes uint64
	AllocMBs   uint64
}

func collectMemoryStatus(title string) memorySnapshot {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	snapshot := memorySnapshot{
		Title:      title,
		AllocBytes: ms.Alloc,
		AllocMBs:   ms.Alloc / (1024 * 1024),
	}

	printMemoryStatus(snapshot)
	return snapshot
}

func printMemoryStatus(snapshot memorySnapshot) {
	data, _ := json.Marshal(snapshot)
	log.Printf("=============%v===========\n", snapshot.Title)
	fmt.Println(string(data))
	log.Println("====================================")
}

func printStats() {
	d, _ := json.Marshal(stats.GetStats())
	log.Println("=============stats==================")
	fmt.Println(string(d))
	log.Println("====================================")

}

func printPCIStats() {
	d, _ := json.Marshal(PCI.FetchStats())
	log.Println("=============pep-coroutine-lib==================")
	fmt.Println(string(d))
	log.Println("====================================")

}
