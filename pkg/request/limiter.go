package request

import (
	"sync"
	"time"
)

func FrequencyLimiter(maxRequests int, periodSeconds time.Duration) func() bool {
	m := sync.Mutex{}
	requests := make([]int64, 0, maxRequests)

	return func() bool {
		m.Lock()
		defer m.Unlock()

		nowUnix := time.Now().Unix()
		newRequests := make([]int64, 0, len(requests))

		for i := range requests {
			diff := nowUnix - requests[i]
			if time.Duration(diff).Seconds() <= periodSeconds.Seconds() {
				newRequests = append(newRequests, requests[i])
			}
		}

		requests = newRequests

		if len(requests) == maxRequests {
			return true
		}

		requests = append(requests, nowUnix)

		return false
	}
}
