package tokenbucket

import (
	"sync/atomic"
	"time"
)

type TokenBucket struct {
	tokens   chan struct{}
	ticker   *time.Ticker
	capacity uint64
	q        chan bool
}

func (t *TokenBucket) Take() bool {
	select {
		case _, ok := <-t.tokens:{ return ok }
		default: return false
	}
}

func (t *TokenBucket) Wait() {
	for !t.Take() {
		time.Sleep(time.Millisecond)
	}
}

func (t *TokenBucket) Close() {
	close(t.q)
}

func (t *TokenBucket) fillTokens() {
	toFillNum := atomic.LoadUint64(&t.capacity) - uint64(len(t.tokens))
	for toFillNum != 0 {
		select {
			case t.tokens <- struct{}{}:{ toFillNum-- }
			default: return
		}
	}
}

func NewTokenBucket(fillInterval time.Duration, capacity uint64) *TokenBucket {
	if fillInterval <= 0 {
		panic("fill interval must > 0")
	}

	if capacity == 0 {
		panic("capacity must > 0")
	}

	tokenBucket := &TokenBucket{tokens: make(chan struct{}, capacity), ticker: time.NewTicker(fillInterval), capacity: capacity, q: make(chan bool)}
	tokenBucket.fillTokens()

	go func(t *TokenBucket) {
		for {
			select {
			case <-t.ticker.C: {
				t.fillTokens()
			}
			case _, available := <-t.q: {
				if !available {
					return
				}
				t.ticker.Stop()
			}
			}
		}
	}(tokenBucket)

	return tokenBucket
}