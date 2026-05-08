package git

import (
	"testing"

	"github.com/cyradin/govnocode/tools/executor"
	"github.com/stretchr/testify/require"
)

func TestCreateBranch_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success: create branch", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		exec := executor.NewShell(tmpDir)

		initRepo(t, tmpDir, exec)

		tool := NewCreateBranch()

		res, err := tool.Execute(t.Context(), exec, []byte(`{"branch":"feature/test"}`))

		require.NoError(t, err)
		require.NotNil(t, res)

		out, err := exec.Execute(t.Context(), []string{"git", "-C", tmpDir, "branch", "--show-current"}, nil)
		require.NoError(t, err)
		require.Equal(t, "feature/test\n", out.Stdout)
	})

	t.Run("error: invalid branch", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		exec := executor.NewShell(tmpDir)

		initRepo(t, tmpDir, exec)

		_, err := NewCreateBranch().Execute(t.Context(), exec, []byte(`{"branch":""}`))

		require.Error(t, err)
	})
}
