package value

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValue_Secret(t *testing.T) {
	duration := 2 * time.Second
	secret := "secret"

	value := New[string](secret, duration, func() {})

	secretFromValue := value.Secret()

	require.EqualValues(t, secretFromValue, secret)
}

func TestValue_TimerStop(t *testing.T) {
	t.Run("Timer Stop Successful", func(t *testing.T) {
		duration := 1 * time.Second
		secret := "secret"

		value := New[string](secret, duration, func() {})

		status := value.TimerStop()

		require.True(t, status)
	})

	t.Run("Timer Stop Unsuccessful", func(t *testing.T) {
		duration := 1 * time.Second
		secret := "secret"

		value := New[string](secret, duration, func() {})
		time.Sleep(2 * time.Second)

		status := value.TimerStop()

		require.False(t, status)
	})
}

func TestValue_TimerReset(t *testing.T) {
	t.Run("Timer Reset Successful", func(t *testing.T) {
		duration := 1 * time.Second
		secret := "secret"

		value := New[string](secret, duration, func() {})

		status := value.TimerReset()

		require.True(t, status)
	})

	t.Run("Timer Reset Unsuccessful", func(t *testing.T) {
		duration := 1 * time.Second
		secret := "secret"

		value := New[string](secret, duration, func() {})
		time.Sleep(2 * time.Second)

		status := value.TimerReset()

		require.False(t, status)
	})
}
