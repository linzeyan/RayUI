package model

// ServerStatItem tracks per-server traffic statistics.
type ServerStatItem struct {
	ProfileID  string `json:"profileId"`
	TotalUp    int64  `json:"totalUp"`
	TotalDown  int64  `json:"totalDown"`
	TodayUp    int64  `json:"todayUp"`
	TodayDown  int64  `json:"todayDown"`
	DateNow    string `json:"dateNow"`
	LastUpdate int64  `json:"lastUpdate"`
}

// TrafficStats is emitted in real time to the frontend.
type TrafficStats struct {
	Up   int64 `json:"up"`   // bytes/sec upload
	Down int64 `json:"down"` // bytes/sec download
}

// SpeedTestResult is the result of a single speed/latency test.
type SpeedTestResult struct {
	ProfileID string `json:"profileId"`
	Latency   int    `json:"latency"`   // milliseconds (-1 = timeout)
	Speed     int64  `json:"speed"`     // bytes/sec (0 if not tested)
}

// CoreStatus describes the running state of a proxy core.
type CoreStatus struct {
	Running   bool      `json:"running"`
	CoreType  ECoreType `json:"coreType"`
	Version   string    `json:"version"`
	StartTime *int64    `json:"startTime,omitempty"`
	PID       int       `json:"pid,omitempty"`
	Profile   string    `json:"profile,omitempty"`
}
