package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyradin/govnocode/internal/command"
	"github.com/stretchr/testify/require"
)

func TestRead_Execute(t *testing.T) {
	t.Parallel()

	t.Run("read file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		file := filepath.Join(tmpDir, "a.txt")
		require.NoError(t, os.WriteFile(file, []byte("hello"), 0644))

		res, err := NewRead().Execute(t.Context(), e, []byte(`{"path":"a.txt"}`))

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Contains(t, res.Stdout, "hello")
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		res, err := NewRead().Execute(t.Context(), e, []byte(`{"path":"qwerty"}`))
		require.Contains(t, res.Stderr, "No such file or directory")
		require.Error(t, err)
	})

	t.Run("invalid path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		_, err := NewRead().Execute(t.Context(), e, []byte(`{"path":""}`))
		require.Error(t, err)
	})
}
