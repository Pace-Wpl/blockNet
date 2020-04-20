package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"time"
)

func NewUUID() (string, error) {
	uuid := make([]byte, 8)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[4] = uuid[4]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[3] = uuid[3]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x", uuid[0:4], uuid[4:8]), nil
}

func GetCurrentTimestampSec() int {
	ts, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))
	return ts
}
