package goerror

import (
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	var stdErr error = errors.New("test_error")

	t.Run("test message", func(t *testing.T) {
		err := New("test_code", "test_message")
		require.Equal(t, "test_code", err.Code)
		require.Equal(t, "test_message", err.Message)

		err = err.WithMessage("test_message_2").(traceableErrorImpl)
		require.Equal(t, "test_message_2", err.Message)
	})

	t.Run("test error wrapping", func(t *testing.T) {
		err := New("test_code", "test_message")
		err = err.WithError(stdErr).(traceableErrorImpl)
		require.Equal(t, stdErr, err.Err)
		require.True(t, errors.Is(err, stdErr))
	})

	t.Run("test origin", func(t *testing.T) {
		err := New("test_code", "test_message")

		_, file, line, ok := runtime.Caller(0)
		require.True(t, ok)
		expectedOrigin := fmt.Sprintf("%s:%d", file, line+3)
		err = err.WithOrigin().(traceableErrorImpl)

		require.Equal(t, expectedOrigin, err.Origin)
	})

	t.Run("test errors.Is", func(t *testing.T) {
		err := New("test_code", "test_message")
		err1 := err.WithError(nil)
		require.True(t, errors.Is(err, err1))

		err2 := New("test_code", "test_message_2")
		require.True(t, errors.Is(err, err2))

		err3 := New("test_code_2", "test_message")
		require.False(t, errors.Is(err, err3))
	})
}
