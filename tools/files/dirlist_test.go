package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyradin/govnocode/internal/command"
	"github.com/stretchr/testify/require"
)

func TestDirList_Execute(t *testing.T) {
	t.Parallel()

	t.Run("list directory", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("x"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("x"), 0644))

		res, err := NewDirList().Execute(t.Context(), e, []byte(`{}`))

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Contains(t, res.Stdout, "a.txt")
		require.Contains(t, res.Stdout, "b.txt")
	})

	t.Run("list subdir", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		sub := filepath.Join(tmpDir, "sub")
		require.NoError(t, os.MkdirAll(sub, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(sub, "c.txt"), []byte("x"), 0644))

		res, err := NewDirList().Execute(t.Context(), e, []byte(`{"path":"sub"}`))

		require.NoError(t, err)
		require.Contains(t, res.Stdout, "c.txt")
	})

	t.Run("empty dir", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		res, err := NewDirList().Execute(t.Context(), e, []byte(`{}`))

		require.NoError(t, err)
		require.Equal(t, "", res.Stderr)
	})
}
