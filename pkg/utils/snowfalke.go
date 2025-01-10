package utils

import (
	"sync"
	"time"
)

const (
	epoch            int64 = 1577836800000 // January 1, 2020, in milliseconds
	machineIDBits    int64 = 10
	dataCenterIDBits int64 = 5
	sequenceBits     int64 = 12
)

var (
	machineID     int64 = 1 // Example machine ID (can be customized)
	dataCenterID  int64 = 1 // Example datacenter ID (can be customized)
	sequence      int64
	lastTimestamp int64 = -1
	mutex         sync.Mutex
)

// GenerateSnowflakeID generates a unique snowflake ID based on timestamp, machine, datacenter, and sequence.
func GenerateSnowflakeID() int64 {
	mutex.Lock()
	defer mutex.Unlock()

	timestamp := time.Now().UnixMilli() - epoch

	if timestamp == lastTimestamp {
		sequence = (sequence + 1) & ((1 << sequenceBits) - 1)
		if sequence == 0 {
			// Wait until the next millisecond
			for timestamp <= lastTimestamp {
				timestamp = time.Now().UnixMilli() - epoch
			}
		}
	} else {
		sequence = 0
	}

	lastTimestamp = timestamp

	// Shift and combine parts
	id := (timestamp << (machineIDBits + dataCenterIDBits + sequenceBits)) |
		(dataCenterID << (machineIDBits + sequenceBits)) |
		(machineID << sequenceBits) |
		sequence

	return id
}
