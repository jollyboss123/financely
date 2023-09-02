package rate

import "time"

type UpdateRequest struct {
	startTime time.Time     `json:"start_time"`
	delay     time.Duration `json:"delay"`
}
