package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyradin/govnocode/internal/command"
	"github.com/stretchr/testify/require"
)

func TestStatus_Execute(t *testing.T) {
	t.Parallel()

	t.Run("empty repo", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		initRepo(t, tmpDir, e)

		res, err := NewStatus().Execute(t.Context(), e, nil)

		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("file appears in status", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		initRepo(t, tmpDir, e)

		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("x"), 0644))

		res, err := NewStatus().Execute(t.Context(), e, nil)

		require.NoError(t, err)
		require.Contains(t, res.Stdout, "a.txt")
	})
}
