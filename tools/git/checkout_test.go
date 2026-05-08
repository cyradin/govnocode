package git

import (
	"testing"

	"github.com/cyradin/govnocode/tools/executor"
	"github.com/stretchr/testify/require"
)

func TestCheckout_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success: switch branch", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		initRepo(t, tmpDir, e)

		_, err := e.Execute(t.Context(), []string{"git", "-C", tmpDir, "branch", "feature/test"}, nil)
		require.NoError(t, err)

		res, err := NewCheckout().Execute(t.Context(), e, []byte(`{"branch":"feature/test"}`))

		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("error: non-existent branch", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		initRepo(t, tmpDir, e)

		_, err := NewCheckout().Execute(t.Context(), e, []byte(`{"branch":"nope"}`))

		require.Error(t, err)
	})
}
