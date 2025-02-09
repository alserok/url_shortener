package utils

import (
	"context"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context) bool
}

const (
	defaultCap  = 10_000
	defaultTick = time.Second
)

func NewLimiter(cap uint, tick time.Duration) *leakyBucket {
	lb := &leakyBucket{
		cap:  defaultCap,
		tick: defaultTick,
	}

	if cap != 0 {
		lb.cap = cap
	}
	if tick.Milliseconds() != 0 {
		lb.tick = tick
	}

	lb.tickets = make(chan struct{}, lb.cap)
	for i := 0; i < int(lb.cap); i++ {
		lb.tickets <- struct{}{}
	}

	go func() {
		defer func() {
			close(lb.tickets)
		}()

		for range time.Tick(lb.tick) {
			select {
			case lb.tickets <- struct{}{}:
			default:
				return
			}
		}
	}()

	return lb
}

type leakyBucket struct {
	cap  uint
	tick time.Duration

	tickets chan struct{}
}

func (l *leakyBucket) Allow(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case <-l.tickets:
		return true
	}
}
