package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyradin/govnocode/tools/executor"
	"github.com/stretchr/testify/require"
)

func TestCommit_Execute(t *testing.T) {
	t.Parallel()

	t.Run("commit file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		initRepo(t, tmpDir, e)

		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("x"), 0644))
		_, _ = e.Execute(t.Context(), []string{"git", "-C", tmpDir, "add", "."}, nil)

		res, err := NewCommit().Execute(t.Context(), e, []byte(`{"message":"test"}`))

		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("empty message", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		initRepo(t, tmpDir, e)

		_, err := NewCommit().Execute(t.Context(), e, []byte(`{"message":""}`))

		require.Error(t, err)
	})
}
