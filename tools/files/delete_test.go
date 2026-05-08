package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyradin/govnocode/internal/command"
	"github.com/stretchr/testify/require"
)

func TestDelete_Execute(t *testing.T) {
	t.Parallel()

	t.Run("delete existing file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		path := filepath.Join(tmpDir, "a.txt")

		require.NoError(
			t,
			os.WriteFile(path, []byte("hello"), 0644),
		)

		res, err := NewDelete().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"a.txt"}`),
		)

		require.NoError(t, err)
		require.NotNil(t, res)

		_, err = os.Stat(path)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("delete missing file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		res, err := NewDelete().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"missing.txt"}`),
		)

		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("empty path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		_, err := NewDelete().Execute(
			t.Context(),
			e,
			[]byte(`{"path":""}`),
		)

		require.Error(t, err)
	})

	t.Run("delete directory recursively", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		dir := filepath.Join(tmpDir, "dir")

		require.NoError(t, os.MkdirAll(filepath.Join(dir, "nested"), 0755))
		require.NoError(
			t,
			os.WriteFile(filepath.Join(dir, "nested", "a.txt"), []byte("x"), 0644),
		)

		res, err := NewDelete().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"dir"}`),
		)

		require.NoError(t, err)
		require.NotNil(t, res)

		_, err = os.Stat(dir)
		require.True(t, os.IsNotExist(err))
	})
}
