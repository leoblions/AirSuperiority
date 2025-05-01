package main

import (
	"time"
)

func CreateDelayToggle(milliseconds int64) func() bool {
	var delayMS = milliseconds
	var startTime = time.Now().UnixMilli()

	var checkTimeExpired = func() bool {
		var now = time.Now().UnixMilli()
		if (now - startTime) > delayMS {
			startTime = now
			return true
		} else {
			return false
		}

	}

	return checkTimeExpired
}
