package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyradin/govnocode/tools/executor"
	"github.com/stretchr/testify/require"
)

func TestMove_Execute(t *testing.T) {
	t.Parallel()

	t.Run("move file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		oldPath := filepath.Join(tmpDir, "old.txt")

		require.NoError(
			t,
			os.WriteFile(oldPath, []byte("hello"), 0644),
		)

		res, err := NewMove().Execute(
			t.Context(),
			e,
			[]byte(`{"from":"old.txt","to":"new.txt"}`),
		)

		require.NoError(t, err)
		require.NotNil(t, res)

		_, err = os.Stat(filepath.Join(tmpDir, "old.txt"))
		require.True(t, os.IsNotExist(err))

		b, err := os.ReadFile(filepath.Join(tmpDir, "new.txt"))
		require.NoError(t, err)

		require.Equal(t, "hello", string(b))
	})

	t.Run("rename file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		require.NoError(
			t,
			os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("x"), 0644),
		)

		_, err := NewMove().Execute(
			t.Context(),
			e,
			[]byte(`{"from":"a.txt","to":"b.txt"}`),
		)

		require.NoError(t, err)

		_, err = os.Stat(filepath.Join(tmpDir, "b.txt"))
		require.NoError(t, err)
	})

	t.Run("source does not exist", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		res, err := NewMove().Execute(
			t.Context(),
			e,
			[]byte(`{"from":"missing.txt","to":"new.txt"}`),
		)

		require.Error(t, err)
		require.NotNil(t, res)
		require.Contains(t, res.Stderr, "No such file or directory")
	})

	t.Run("empty from", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		_, err := NewMove().Execute(
			t.Context(),
			e,
			[]byte(`{"from":"","to":"b.txt"}`),
		)

		require.Error(t, err)
	})

	t.Run("empty to", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		_, err := NewMove().Execute(
			t.Context(),
			e,
			[]byte(`{"from":"a.txt","to":""}`),
		)

		require.Error(t, err)
	})
}
