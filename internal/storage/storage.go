package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
)

type Record struct {
	Key   string
	Value string

	Mu sync.Mutex
}

type KVStore interface {
	Set(r Record) error
	Get(key string) (*Record, error)
}

type kvStore struct {
	Storage  []*Record
	Capacity int
	Locks    []sync.Mutex
}

func NewKVStore(capacity int) KVStore {
	return &kvStore{
		Capacity: capacity,
		Storage:  make([]*Record, capacity),
		Locks:    make([]sync.Mutex, capacity),
	}
}

func (k *kvStore) Set(r Record) error {
	hash := formatKey(r.Key, int64(k.Capacity))

	if k.Storage[hash] != nil && r.Key != k.Storage[hash].Key {
		// Hash collision
		return fmt.Errorf("hash collisioned")
	}

	// Update or Insert
	k.Locks[hash].Lock()
	k.Storage[hash] = &r
	k.Locks[hash].Unlock()

	return nil
}

func (k *kvStore) Get(key string) (*Record, error) {
	hash := formatKey(key, int64(k.Capacity))

	if k.Storage[hash] == nil {
		return nil, fmt.Errorf("the key is not found")
	}

	return k.Storage[hash], nil
}

func formatKey(key string, capacity int64) int {
	bi := big.NewInt(0)
	h := md5.Sum([]byte(key))
	hexStr := hex.EncodeToString(h[:])
	bi.SetString(hexStr, 16)
	decimal := bi.Int64()

	hash := decimal % capacity

	if hash < 0 {
		hash += capacity
	}

	return int(hash)
}
