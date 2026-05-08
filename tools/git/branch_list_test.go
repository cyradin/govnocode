package git

import (
	"testing"

	"github.com/cyradin/govnocode/internal/command"
	"github.com/stretchr/testify/require"
)

func TestBranchList_Execute(t *testing.T) {
	t.Parallel()

	t.Run("list branches", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		initRepo(t, tmpDir, e)

		_, _ = e.Execute(t.Context(), []string{"git", "-C", tmpDir, "branch", "feature/test"}, nil)

		res, err := NewBranchList().Execute(t.Context(), e, nil)

		require.NoError(t, err)
		require.Contains(t, res.Stdout, "test")
	})
}
