package db


import (
	"hash/fnv"
	"sort"
)
func HashKey(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return h.Sum64()
}

type KeyIndex struct {
	hashes  []uint64
	offsets []int64
	count   int
}

func NewKeyIndex(capacity int) *KeyIndex {
	return &KeyIndex{
		hashes:  make([]uint64, capacity),
		offsets: make([]int64, capacity),
		count:   0,
	}
}

func (ki *KeyIndex) Add(keyHash uint64, offset int64) {
	if ki.count < len(ki.hashes) {
		ki.hashes[ki.count] = keyHash
		ki.offsets[ki.count] = offset
		ki.count++
	}
}

func (ki *KeyIndex) Sort() {
	type kv struct {
		hash   uint64
		offset int64
	}

	// Create view to sort
	kvs := make([]kv, ki.count)
	for i := 0; i < ki.count; i++ {
		kvs[i] = kv{ki.hashes[i], ki.offsets[i]}
	}

	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].hash < kvs[j].hash
	})

	for i := 0; i < ki.count; i++ {
		ki.hashes[i] = kvs[i].hash
		ki.offsets[i] = kvs[i].offset
	}
}

func (ki *KeyIndex) Get(keyHash uint64) (int64, bool) {
	// Binary search
	idx := sort.Search(ki.count, func(i int) bool {
		return ki.hashes[i] >= keyHash
	})

	if idx < ki.count && ki.hashes[idx] == keyHash {
		return ki.offsets[idx], true
	}
	return 0, false
}