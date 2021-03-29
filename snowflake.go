package common

import (
	"sync"
	"time"
)

const (
	snowflakeMaxSerialNumber uint16 = uint16(0b00001111)<<8 + 0xff
	snowflake63Bit                  = ^(uint64(1) << 63)
)

var (
	snowflakeMutex        sync.Mutex
	snowflakeTimestamp    int64
	snowflakeSerialNumber uint16
)

func SnowflakeID(groupID, mechineID byte) (n uint64) {
	// 时间戳
	timestamp := time.Now().UTC().Unix()
	// 判断在这个毫秒内生成的个数
	var serialNumber uint16
	snowflakeMutex.Lock()
	if timestamp == snowflakeTimestamp {
		snowflakeSerialNumber++
		// 1ms内生成的个数太多了，加1ms
		if snowflakeSerialNumber > snowflakeMaxSerialNumber {
			snowflakeSerialNumber = 0
			timestamp++
			snowflakeTimestamp = timestamp
		}
	}
	serialNumber = snowflakeSerialNumber
	snowflakeMutex.Unlock()
	// 12bit id serial number
	n |= uint64(serialNumber)
	// 5bit mechine id
	n |= uint64(mechineID&0b0011111) << 12
	// 5bit group id
	n |= uint64(groupID&0b0011111) << 17
	// 41bit timestamp
	n |= uint64(timestamp) << 22
	// 1bit 0
	n &= snowflake63Bit
	return n
}
