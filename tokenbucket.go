package tokenbucket

import (
	"sync/atomic"
	"time"
)

type TokenBucket struct {
	tokens   chan struct{}
	ticker   *time.Ticker
	capacity uint32
}

func (t *TokenBucket) SetCapacity(num uint32) {
	atomic.StoreUint32(&t.capacity, num)
}

func (t *TokenBucket) Capacity() uint32 {
	return atomic.LoadUint32(&t.capacity)
}

func (t *TokenBucket) Take() int {
	select {
		case _, ok := <-t.tokens: if ok {return len(t.tokens)}; return 0
		default: return 0
	}
}

func (t *TokenBucket) Wait() int {
	<- t.tokens
	return len(t.tokens)
}

func (t *TokenBucket) Close() {
	t.ticker.Stop()
}

func (t *TokenBucket) fillTokens() {
	toFillNum := atomic.LoadUint32(&t.capacity) - uint32(len(t.tokens))
	for toFillNum != 0 {
		select {
			case t.tokens <- struct{}{}: toFillNum--
			default: return
		}
	}
}

func NewTokenBucket(fillInterval time.Duration, capacity uint32) *TokenBucket {
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