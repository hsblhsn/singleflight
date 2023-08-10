// Package singleflight provides a duplicate function call suppression.
// It's a generic supported wrapper of golang.org/x/sync/singleflight.
package singleflight

import "golang.org/x/sync/singleflight"

// Group represents a class of work and forms a namespace in
// which units of work can be executed with duplicate suppression.
type Group[T any] struct {
	group *singleflight.Group
}

// Result holds the results of Do, so they can be passed
// on a channel.
type Result[T any] struct {
	Val    T
	Err    error
	Shared bool
}

// NewGroup creates a new Group instance.
func NewGroup[T any]() *Group[T] {
	return &Group[T]{
		group: &singleflight.Group{},
	}
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
// The return value shared indicates whether v was given to multiple callers.
func (g *Group[T]) Do(key string, fn func() (T, error)) (T, error, bool) {
	v, err, shared := g.group.Do(key, func() (interface{}, error) {
		return fn()
	})
	return v.(T), err, shared
}

// DoChan is like Do but returns a channel that will receive the
// results when they are ready.
//
// The returned channel will not be closed.
func (g *Group[T]) DoChan(key string, fn func() (T, error)) <-chan Result[T] {
	ch := make(chan Result[T], 1)
	go func() {
		v, err, shared := g.group.Do(key, func() (interface{}, error) {
			return fn()
		})
		ch <- Result[T]{
			Val:    v.(T),
			Err:    err,
			Shared: shared,
		}
	}()
	return ch
}

// Forget tells the singleflight to forget about a key.  Future calls
// to Do for this key will call the function rather than waiting for
// an earlier call to complete.
func (g *Group[T]) Forget(key string) {
	g.group.Forget(key)
}
