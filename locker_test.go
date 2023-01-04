package locker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLocker_Status(t *testing.T) {
	locker := New[string, string]()

	t.Run("Status Failure", func(t *testing.T) {
		key := "key"
		status := locker.Status(key)

		require.False(t, status)
	})

	t.Run("Status Success", func(t *testing.T) {
		key := "key"
		secret := "secret"
		duration := 1 * time.Second

		locker.Lock(key, secret, duration, nil)

		status := locker.Status(key)

		require.True(t, status)
	})
}

func TestLocker_Lock(t *testing.T) {
	locker := New[string, string]()

	t.Run("Lock Success", func(t *testing.T) {
		key := "key"
		secret := "secret"
		duration := 1 * time.Second

		err := locker.Lock(key, secret, duration, nil)

		require.NoError(t, err)
	})

	t.Run("Lock Failure", func(t *testing.T) {
		key := "key"
		secret := "secret"
		duration := 1 * time.Second

		err := locker.Lock(key, secret, duration, nil)

		require.ErrorIs(t, err, ErrLockExists)
	})

}

func TestLocker_Unlock(t *testing.T) {
	locker := New[string, string]()

	t.Run("Failure: Key does not exist", func(t *testing.T) {
		key := "key_2"
		secret := "secret_2"

		_, err := locker.Unlock(key, secret)

		require.ErrorIs(t, err, ErrInvalidKey)
	})

	t.Run("Failure: secret mismatch", func(t *testing.T) {
		key := "key_1"
		secret := "secret_1"
		duration := 1 * time.Second
		locker.Lock(key, secret, duration, nil)

		secret = "secret_2"
		_, err := locker.Unlock(key, secret)

		require.ErrorIs(t, err, ErrInvalidSecret)
	})

	t.Run("Success: unlocked", func(t *testing.T) {
		key := "key_1"
		secret := "secret_1"
		duration := 1 * time.Second
		locker.Lock(key, secret, duration, nil)

		status, err := locker.Unlock(key, secret)

		require.NoError(t, err)
		require.True(t, status)
	})
}
