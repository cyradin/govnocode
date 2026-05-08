package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyradin/govnocode/internal/command"
	"github.com/stretchr/testify/require"
)

func TestMkdir_Execute(t *testing.T) {
	t.Parallel()

	t.Run("create directory", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		res, err := NewMkdir().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"testdir"}`),
		)

		require.NoError(t, err)
		require.NotNil(t, res)

		info, err := os.Stat(filepath.Join(tmpDir, "testdir"))
		require.NoError(t, err)
		require.True(t, info.IsDir())
	})

	t.Run("create nested directories", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		_, err := NewMkdir().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"a/b/c"}`),
		)

		require.NoError(t, err)

		info, err := os.Stat(filepath.Join(tmpDir, "a/b/c"))
		require.NoError(t, err)
		require.True(t, info.IsDir())
	})

	t.Run("empty path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		_, err := NewMkdir().Execute(
			t.Context(),
			e,
			[]byte(`{"path":""}`),
		)

		require.Error(t, err)
	})

	t.Run("mkdir already exists", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		require.NoError(
			t,
			os.Mkdir(filepath.Join(tmpDir, "testdir"), 0755),
		)

		_, err := NewMkdir().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"testdir"}`),
		)

		require.NoError(t, err)
	})
}
