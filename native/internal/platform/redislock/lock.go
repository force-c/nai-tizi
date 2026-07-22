// Package redislock provides a small context-aware Redis distributed lock.
package redislock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	randv2 "math/rand/v2"
	"time"

	"github.com/redis/go-redis/v9"
)

const releaseScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
end
return 0
`

var ErrNotOwner = errors.New("redis lock is no longer owned by this holder")

// Locker creates independent locks backed by one Redis client.
type Locker struct {
	client   *redis.Client
	ttl      time.Duration
	retryMin time.Duration
	retryMax time.Duration
}

// Option customizes a Locker.
type Option func(*Locker)

// WithTTL sets the lock expiry.
func WithTTL(ttl time.Duration) Option {
	return func(l *Locker) { l.ttl = ttl }
}

// WithRetryRange sets the randomized acquisition retry interval.
func WithRetryRange(minimum, maximum time.Duration) Option {
	return func(l *Locker) {
		l.retryMin = minimum
		l.retryMax = maximum
	}
}

// New creates a Locker.
func New(client *redis.Client, options ...Option) *Locker {
	l := &Locker{
		client:   client,
		ttl:      10 * time.Second,
		retryMin: 20 * time.Millisecond,
		retryMax: 80 * time.Millisecond,
	}
	for _, option := range options {
		option(l)
	}
	if l.retryMax < l.retryMin {
		l.retryMax = l.retryMin
	}
	return l
}

// Lock is one acquired lock instance.
type Lock struct {
	client *redis.Client
	key    string
	token  string
}

// Acquire waits until the lock is acquired or ctx is cancelled.
func (l *Locker) Acquire(ctx context.Context, key string) (*Lock, error) {
	if l == nil || l.client == nil {
		return nil, errors.New("redis lock client is nil")
	}
	if key == "" {
		return nil, errors.New("redis lock key is empty")
	}
	if l.ttl <= 0 {
		return nil, errors.New("redis lock ttl must be positive")
	}
	token, err := newToken()
	if err != nil {
		return nil, fmt.Errorf("create redis lock token: %w", err)
	}
	for {
		acquired, err := l.client.SetNX(ctx, key, token, l.ttl).Result()
		if err != nil {
			return nil, fmt.Errorf("acquire redis lock %q: %w", key, err)
		}
		if acquired {
			return &Lock{client: l.client, key: key, token: token}, nil
		}
		timer := time.NewTimer(l.retryDelay())
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, fmt.Errorf("acquire redis lock %q: %w", key, ctx.Err())
		case <-timer.C:
		}
	}
}

// Release removes the lock only if its random token is still the owner.
func (l *Lock) Release(ctx context.Context) error {
	if l == nil || l.client == nil {
		return errors.New("redis lock is nil")
	}
	result, err := l.client.Eval(ctx, releaseScript, []string{l.key}, l.token).Int64()
	if err != nil {
		return fmt.Errorf("release redis lock %q: %w", l.key, err)
	}
	if result == 0 {
		return ErrNotOwner
	}
	return nil
}

// WithLock acquires key, runs fn, and safely releases the lock.
func (l *Locker) WithLock(ctx context.Context, key string, fn func() error) (err error) {
	lock, err := l.Acquire(ctx, key)
	if err != nil {
		return err
	}
	defer func() {
		releaseCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		releaseErr := lock.Release(releaseCtx)
		if releaseErr != nil && !errors.Is(releaseErr, ErrNotOwner) {
			err = errors.Join(err, releaseErr)
		}
	}()
	return fn()
}

func (l *Locker) retryDelay() time.Duration {
	if l.retryMax <= l.retryMin {
		return l.retryMin
	}
	return l.retryMin + time.Duration(randv2.Int64N(int64(l.retryMax-l.retryMin)+1))
}

func newToken() (string, error) {
	var data [16]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(data[:]), nil
}
