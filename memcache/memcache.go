// Package memcache contains helpers for caching function results.
package memcache

import (
	"cmp"
	"context"
	"log/slog"
	"sync"
	"time"
)

// Cache implements a cache.
type Cache[T any] interface {

	// Get gets the current value, updating it if necessary. The returned
	// pointer is shared and should not be modified. If the value has not been
	// updated, the returned pointer must be the same. If err is nil, the
	// returned pointer must not be nil. If an update failed but the cached data
	// is still valid, both a pointer and an error may be returned.
	Get() (*T, error)
}

// CacheFunc wraps a func implementing [Cache].
type CacheFunc[T any] func() (*T, error)

func (fn CacheFunc[T]) Get() (*T, error) {
	v, err := fn()
	if err == nil && v == nil {
		panic("cache must return a non-nil pointer if err is nil")
	}
	return v, err
}

// MultiCache dynamically initializes caches.
func MultiCache[K comparable, T any](init func(K) Cache[T]) func(K) Cache[T] {
	var (
		cacheMu  sync.RWMutex
		cacheMap = map[K]Cache[T]{}
	)
	return func(k K) Cache[T] {
		cacheMu.RLock()
		c := cacheMap[k]
		cacheMu.RUnlock()
		if c == nil {
			cacheMu.Lock()
			if c = cacheMap[k]; c == nil {
				c = init(k)
				cacheMap[k] = c
			}
			cacheMu.Unlock()
		}
		return c
	}
}

// Backoff implements a backoff strategy.
type Backoff interface {

	// Backoff returns the next time to retry after the specified error at the
	// specified time on the specified attempt (>= 1).
	Backoff(t time.Time, err error, attempt int) time.Time
}

// BackoffFunc wraps a func implementing [Backoff].
type BackoffFunc func(t time.Time, err error, attempt int) time.Time

func (fn BackoffFunc) Backoff(t time.Time, err error, attempt int) time.Time {
	return fn(t, err, attempt)
}

// CacheConfig configures [Cache].
type CacheConfig struct {

	// Timeout is the maximum amount of time an update attempt can block for. If
	// negative, no timeout is used. If zero, the default value is used.
	Timeout time.Duration

	// CacheTime is the maximum amount of time data is cached for before
	// attempting to update it. If negative, an update is attempted every time.
	// If zero, the default value is used.
	CacheTime time.Duration

	// StaleTime is the maximum amount of time to return old data (along with
	// the update error) while updates fail. If negative, old data will never be
	// returned. If zero, the default value is used.
	StaleTime time.Duration

	// Backoff is used to delay update retries on error. If nil, no backoff is
	// used.
	Backoff Backoff

	// Logger is used to write informational logs about cache updates. If nil,
	// no logger is used.
	Logger *slog.Logger
}

// Cached wraps the provided fetch function in a cache.
func Cached[T any](cfg CacheConfig, fetch func(ctx context.Context) (T, error)) Cache[T] {
	cfg.Timeout = negZeroDef(cfg.Timeout, time.Second*7)
	cfg.CacheTime = negZeroDef(cfg.CacheTime, time.Minute*15)
	cfg.StaleTime = negZeroDef(cfg.StaleTime, time.Hour*2)

	var cache struct {
		mu sync.Mutex

		failure  time.Time
		failureV error
		failureN int

		success  time.Time
		successV *T
	}
	if cfg.Logger != nil {
		cfg.Logger.Info("cache created", slog.Group("config", "timeout", cfg.Timeout.Seconds(), "cache_time", cfg.CacheTime.Seconds(), "stale_time", cfg.StaleTime.Seconds(), "backoff", cfg.Backoff != nil))
	}
	return CacheFunc[T](func() (*T, error) {
		cache.mu.Lock()
		defer cache.mu.Unlock()

		ctx := context.Background()

		if cfg.Timeout > 0 {
			var cancel func()
			ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
			defer cancel()
		}

		now := time.Now()

		if !cache.success.IsZero() {
			age := time.Since(cache.success)
			if age <= cfg.CacheTime {
				if cfg.Logger != nil {
					cfg.Logger.Debug("using cached data", "age", age.Truncate(time.Millisecond).Seconds())
				}
				return cache.successV, nil
			}
			if age > cfg.CacheTime+cfg.StaleTime {
				if cfg.Logger != nil {
					cfg.Logger.Debug("clearing stale cached data", "age", age.Truncate(time.Millisecond).Seconds())
				}
				cache.success = time.Time{}
				cache.successV = nil
			}
		}

		if cfg.Backoff != nil {
			if cache.failureN != 0 {
				if t := cfg.Backoff.Backoff(cache.failure, cache.failureV, cache.failureN); !t.IsZero() && now.Before(t) {
					if cfg.Logger != nil {
						if cache.success.IsZero() {
							cfg.Logger.Debug("no cached data to use")
						} else {
							cfg.Logger.Debug("using old cached data")
						}
						cfg.Logger.Debug("not updating cached data due to backoff", "attempt", cache.failureN, "error", cache.failureV, "error_at", cache.failure, "backoff_until", t)
					}
					return cache.successV, cache.failureV
				}
			}
		}

		if cfg.Logger != nil {
			cfg.Logger.Info("updating cached data", "attempt", cache.failureN)
		}

		if v, err := forceContextCancel1(ctx, fetch); err != nil {
			cache.failure = now
			cache.failureV = err
			cache.failureN++
		} else {
			cache.failure = time.Time{}
			cache.failureV = nil
			cache.failureN = 0
			cache.success = now
			cache.successV = &v
		}
		if cfg.Logger != nil {
			if !cache.failure.IsZero() {
				if cfg.Backoff != nil {
					cfg.Logger.Warn("failed to update cached data", "attempt", cache.failureN, "duration", time.Since(now).Truncate(time.Millisecond).Seconds(), "error", cache.failureV, "backoff", cfg.Backoff.Backoff(cache.failure, cache.failureV, cache.failureN), "using_old_data", !cache.success.IsZero())
				} else {
					cfg.Logger.Warn("failed to update cached data", "attempt", cache.failureN, "duration", time.Since(now).Truncate(time.Millisecond).Seconds(), "error", cache.failureV, "using_old_data", !cache.success.IsZero())
				}
				if cache.success.IsZero() {
					cfg.Logger.Debug("no cached data to use")
				} else {
					cfg.Logger.Debug("using old cached data")
				}
			} else {
				cfg.Logger.Info("successfully updated cached data", "attempt", cache.failureN, "duration", time.Since(now).Truncate(time.Millisecond).Seconds())
			}
		}
		return cache.successV, cache.failureV
	})
}

// CachedTransformConfig configures [CachedTransform].
type CachedTransformConfig struct {

	// Logger is used to write informational logs about cache updates. If nil,
	// no logger is used.
	Logger *slog.Logger
}

// CachedTransform transforms the value from a cache, updating it only when it
// changes, or if the source returns an update error. Note that unlike [Cached],
// if the function errors, only an error is returned, and otherwise, only a
// value is returned.
func CachedTransform[T, U any](source Cache[T], cfg CachedTransformConfig, transform func(v T, err error) (U, error)) Cache[U] {
	var cache struct {
		mu     sync.Mutex
		src    *T
		srcErr error
		res    *U
		resErr error
	}
	if cfg.Logger != nil {
		cfg.Logger.Info("cache transform created")
	}
	return CacheFunc[U](func() (*U, error) {
		cache.mu.Lock()
		defer cache.mu.Unlock()

		src, srcErr := source.Get()
		if srcErr != nil {
			if src == nil {
				if cfg.Logger != nil {
					cfg.Logger.Debug("clearing transform result since transform source update failed and doesn't have old cached data")
				}
				cache.res = nil
				cache.resErr = nil
				return nil, srcErr
			}
		} else if src == nil {
			panic("cache must return a non-nil pointer if err is nil")
		}

		if cache.src != src || cache.srcErr != srcErr {
			now := time.Now()

			if cfg.Logger != nil {
				cfg.Logger.Info("executing transform")
			}

			res, resErr := transform(*src, srcErr)
			cache.src, cache.srcErr = src, srcErr // note: it's important that this is after transform so if it panics, it will try again
			cache.res, cache.resErr = &res, resErr

			if cfg.Logger != nil {
				if cache.resErr != nil {
					cfg.Logger.Warn("failed to execute transform", "duration", time.Since(now).Truncate(time.Millisecond).Seconds(), "error", cache.resErr)
				} else {
					cfg.Logger.Info("successfully executed transform", "duration", time.Since(now).Truncate(time.Millisecond).Seconds())
				}
			}
		}
		return cache.res, cache.resErr
	})
}

// forceContextCancel runs fn, immediately returning on context cancellation
// (leaving fn running in the background until it exits on its own).
func forceContextCancel(ctx context.Context, fn func(context.Context) error) error {
	var (
		ret = make(chan struct{})
		err error
	)
	go func() {
		defer close(ret)
		err = fn(ctx)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ret:
		return err
	}
}

// forceContextCancel1 wraps forceContextCancel for a single return value.
func forceContextCancel1[T any](ctx context.Context, fn func(context.Context) (T, error)) (T, error) {
	var ret1 T
	err := forceContextCancel(ctx, func(ctx context.Context) (err error) {
		ret1, err = fn(ctx)
		return err
	})
	return ret1, err
}

// negZerDef returns def if val is zero, zero if val is negative, and val
// otherwise.
func negZeroDef[T cmp.Ordered](val, def T) T {
	var zero T
	switch {
	case val == zero:
		return def
	case val < zero:
		return zero
	default:
		return val
	}
}
