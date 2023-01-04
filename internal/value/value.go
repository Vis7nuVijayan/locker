package value

import "time"

type Value[T comparable] struct {
	secret T

	timer    *time.Timer
	duration time.Duration
	onExpiry func()
}

// New creates a struct that stores a secret and a timer for expiry of the struct
func New[T comparable](secret T, duration time.Duration, onExpiry func()) *Value[T] {
	return &Value[T]{
		secret:   secret,
		timer:    time.AfterFunc(duration, onExpiry),
		duration: duration,
		onExpiry: onExpiry,
	}
}

// Secret returns secret stored in Value struct
func (v Value[T]) Secret() T {
	return v.secret
}

// TimerStop stops the running timer
// returns true is the stop is successful
// returns false if the timer has already expired
func (v Value[T]) TimerStop() bool {
	return v.timer.Stop()
}

// TimerReset resets the running timer
// returns true is the reset is successful
// returns false if the timer has already expired
func (v *Value[T]) TimerReset() bool {
	if !v.timer.Stop() {
		return false
	}

	v.timer = time.AfterFunc(v.duration, v.onExpiry)

	return true
}
