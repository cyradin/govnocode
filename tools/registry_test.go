package tools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistry_Register(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		r := NewRegistry()

		err := r.Register(&Tool{spec: Spec{Code: "code"}})
		require.NoError(t, err)
	})

	t.Run("success, multiple tools", func(t *testing.T) {
		t.Parallel()

		r := NewRegistry()

		err := r.Register(
			&Tool{spec: Spec{Code: "code1"}},
			&Tool{spec: Spec{Code: "code2"}},
		)
		require.NoError(t, err)
	})

	t.Run("duplicate error", func(t *testing.T) {
		t.Parallel()

		r := NewRegistry()

		err := r.Register(&Tool{spec: Spec{Code: "code"}})
		require.NoError(t, err)

		err = r.Register(&Tool{spec: Spec{Code: "code"}})
		require.ErrorIs(t, err, ErrAlreadyRegistered)
	})
}

func TestRegistry_Get(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tool := &Tool{spec: Spec{Code: "code"}}

		r := NewRegistry()

		err := r.Register(tool)
		require.NoError(t, err)

		result, err := r.Get(tool.spec.Code)
		require.NoError(t, err)

		require.Equal(t, tool, result)
	})

	t.Run("not found error", func(t *testing.T) {
		t.Parallel()

		r := NewRegistry()

		result, err := r.Get("code")

		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, result)
	})
}
