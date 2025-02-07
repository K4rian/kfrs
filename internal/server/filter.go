package server

import "time"

func filterRecentRequests(times []time.Time, currentTime time.Time) []time.Time {
	var recent []time.Time
	for _, t := range times {
		if currentTime.Sub(t) <= time.Minute {
			recent = append(recent, t)
		}
	}
	return recent
}
