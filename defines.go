package main

import (
  "time"
)

const (
  maxRetryCount = 5
  retryBufferCapability = 10
  retryBaseDelayTime = 5 * time.Second
  csvHeader = "timestamp,freeway_id,location_id,direction_1,speed_1,direction_2,speed_2\r\n"
)

type retryInfo struct {
  date time.Time
  delayTime time.Duration
  retryCount int
}
