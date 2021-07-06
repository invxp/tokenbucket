package tokenbucket

import (
	"sync/atomic"
	"time"
)

type TokenBucket struct {
	tokens   chan struct{}
	ticker   *time.Ticker
	capacity uint64
}

func (t *TokenBucket) Take() bool {
	select {
		case _, ok := <-t.tokens: return ok
		default: return false
	}
}

func (t *TokenBucket) Wait() {
	<- t.tokens
}

func (t *TokenBucket) Close() {
	t.ticker.Stop()
}

func (t *TokenBucket) fillTokens() {
	toFillNum := atomic.LoadUint64(&t.capacity) - uint64(len(t.tokens))
	for toFillNum != 0 {
		select {
			case t.tokens <- struct{}{}: toFillNum--
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

	tokenBucket := &TokenBucket{tokens: make(chan struct{}, capacity), ticker: time.NewTicker(fillInterval), capacity: capacity}
	tokenBucket.fillTokens()

	go func(t *TokenBucket) {
		for range t.ticker.C {
			t.fillTokens()
		}
	}(tokenBucket)

	return tokenBucket
}