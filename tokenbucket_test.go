package tokenbucket

import (
	"testing"
	"time"
)

func TestTokenBucketWait(t *testing.T) {
	tokenBucket := NewTokenBucket(time.Second, 1000)

	for {
		last := time.Now()
		tokenBucket.Wait()
		t.Log("WAIT:", time.Since(last))
	}
}

func TestTokenBucketTake(t *testing.T) {
	tokenBucket := NewTokenBucket(time.Second, 1000)

	for {
		t.Log("TAKE:", tokenBucket.Take())
	}
}