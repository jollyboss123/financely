package cron

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type Job struct {
	ID        string
	cancel    context.CancelFunc
	startTime time.Time
	delay     time.Duration
}

type JobFunc func(time.Time)

var jobs = make(map[string]*Job)
var mu sync.Mutex

func Start(jobID string, startTime time.Time, delay time.Duration, jobFunc JobFunc) (string, error) {
	mu.Lock()
	_, exists := jobs[jobID]
	mu.Unlock()

	if exists {
		return "", errors.New("job id already exists")
	}

	ctx, cancel := context.WithCancel(context.Background())

	job := &Job{
		ID:        jobID,
		cancel:    cancel,
		startTime: startTime,
		delay:     delay,
	}

	mu.Lock()
	jobs[jobID] = job
	mu.Unlock()

	go func() {
		for t := range cron(ctx, startTime, delay) {
			jobFunc(t)
			log.Printf("job: %s run at %v\n", jobID, t.Format("2006-01-02 15:04:05"))
		}
	}()

	return jobID, nil
}

func Cancel(jobID string) {
	mu.Lock()
	job, exists := jobs[jobID]
	mu.Unlock()

	if exists {
		job.cancel()
	}
}

func cron(ctx context.Context, startTime time.Time, delay time.Duration) <-chan time.Time {
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
