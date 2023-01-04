package locker

import (
	"errors"
	"sync"
	"time"

	"github.com/Vis7nuVijayan/locker/internal/value"
)

var (
	ErrLockExists    = errors.New("lock already exists")
	ErrInvalidKey    = errors.New("invalid key")
	ErrInvalidSecret = errors.New("invalid secret")
)

type Locker[K comparable, V comparable] struct {
	lock    sync.RWMutex
	lockMap map[K]*value.Value[V]
}

func New[K comparable, V comparable]() *Locker[K, V] {
	return &Locker[K, V]{
		lock:    sync.RWMutex{},
		lockMap: make(map[K]*value.Value[V]),
	}
}

// Status returns the status of the lock for the key
func (l *Locker[K, V]) Status(key K) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	_, exists := l.lockMap[key]

	return exists
}

// Lock create a new entry for locking for the key if no lock exists
func (l *Locker[K, V]) Lock(key K, secret V, duration time.Duration, expire func()) error {
	if l.Status(key) {
		return ErrLockExists
	}

	if expire == nil {
		expire = func() {}
	}

	// onExpiry calls the user provided expire function and also the Remove function to delete the key from locker
	onExpiry := func() {
		expire()

		l.Remove(key)
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	l.lockMap[key] = value.New(secret, duration, onExpiry)

	return nil
}

// Unlock unlocks the lock for the key if the provided secret is the same stored in the struct
func (l *Locker[K, V]) Unlock(key K, secret V) (bool, error) {
	value, err := l.get(key, secret)
	if err != nil {
		return false, err
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	value.TimerStop()

	l.remove(key)

	return true, nil
}

// Reset returns true/false if it is able to reset the timer of the expiration of the lock
func (l *Locker[K, V]) Reset(key K, secret V) (bool, error) {
	value, err := l.get(key, secret)
	if err != nil {
		return false, err
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	return value.TimerReset(), nil
}

// Remove delete value from the locker map
func (l *Locker[K, V]) Remove(key K) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.remove(key)

	return nil
}

func (l *Locker[K, V]) remove(key K) {
	delete(l.lockMap, key)
}

// get return the struct associated with the key if the key exists and the secrets match
func (l *Locker[K, V]) get(key K, secret V) (*value.Value[V], error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	value, exists := l.lockMap[key]
	if !exists {
		return nil, ErrInvalidKey
	}

	if value.Secret() != secret {
		return nil, ErrInvalidSecret
	}

	return value, nil
}
