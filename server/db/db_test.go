package db

import (
	"fmt"
	"runtime"
	"testing"
)

func TestMemoryOverhead(t *testing.T) {
	var m1, m2, m3 runtime.MemStats

	runtime.ReadMemStats(&m1)
	fmt.Printf("Before DbScan: %v MB\n", m1.Alloc/1024/1024)

	index, err := DbScan()
	if err != nil {
		t.Fatal(err)
	}

	runtime.ReadMemStats(&m2)
	fmt.Printf("After DbScan: %v MB\n", m2.Alloc/1024/1024)
	fmt.Printf("Memory used (before GC): %v MB\n", (m2.Alloc-m1.Alloc)/1024/1024)

	// Force GC to clear temporary allocations from sort
	runtime.GC()
	runtime.ReadMemStats(&m3)
	fmt.Printf("After GC: %v MB\n", m3.Alloc/1024/1024)
	fmt.Printf("Memory used (after GC): %v MB\n", (m3.Alloc-m1.Alloc)/1024/1024)

	fmt.Printf("Keys in index: %v\n", index.count)
	fmt.Printf("Bytes per key (after GC): %v\n", (m3.Alloc-m1.Alloc)/uint64(index.count))
}

func BenchmarkDbScan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DbScan()
	}
}
