package cron

import (
	"context"
	"time"
)

func Cron(ctx context.Context, startTime time.Time, delay time.Duration) <-chan time.Time {
	stream := make(chan time.Time, 1)

	// Need to check if the time is zero (e.g. if time.Time{} was used)
	if !startTime.IsZero() {
		diff := time.Until(startTime)
		if diff < 0 {
			total := diff - delay
			times := total / delay * -1

			startTime = startTime.Add(times * delay)
		}
	}

	go func() {
		t := <-time.After(time.Until(startTime))
		stream <- t

		ticker := time.NewTicker(delay)
		defer ticker.Stop()

		for {
			select {
			case t2 := <-ticker.C:
				stream <- t2
			case <-ctx.Done():
				close(stream)
				return
			}
		}
	}()

	return stream
}
