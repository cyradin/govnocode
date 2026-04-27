package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileRead_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success: read file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewFileRead()

		file := filepath.Join(tmpDir, "a.txt")
		require.NoError(t, os.WriteFile(file, []byte("hello"), 0644))

		res, err := tool.Execute(tmpDir, []byte(`{"path":"a.txt"}`))

		require.NoError(t, err)
		require.Contains(t, res.Stdout, "hello")
	})

	t.Run("error: file not exists", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewFileRead()

		_, err := tool.Execute(tmpDir, []byte(`{"path":"no.txt"}`))

		require.Error(t, err)
	})
}

func TestFileWrite_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success: write file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewFileWrite()

		res, err := tool.Execute(tmpDir, []byte(`{"path":"a.txt","content":"hello"}`))

		require.NoError(t, err)
		require.NotNil(t, res)

		b, err := os.ReadFile(filepath.Join(tmpDir, "a.txt"))
		require.NoError(t, err)
		require.Equal(t, "hello", string(b))
	})

	t.Run("error: empty path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewFileWrite()

		_, err := tool.Execute(tmpDir, []byte(`{"path":"","content":"x"}`))

		require.Error(t, err)
	})
}

func TestFileDelete_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success: delete file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewFileDelete()

		file := filepath.Join(tmpDir, "a.txt")
		require.NoError(t, os.WriteFile(file, []byte("x"), 0644))

		res, err := tool.Execute(tmpDir, []byte(`{"path":"a.txt"}`))

		require.NoError(t, err)
		require.NotNil(t, res)

		_, err = os.Stat(file)
		require.Error(t, err)
	})

	t.Run("error: missing file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewFileDelete()

		_, err := tool.Execute(tmpDir, []byte(`{"path":"no.txt"}`))

		require.Error(t, err)
	})
}

func TestDirList_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success: list directory", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewDirList()

		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("x"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("x"), 0644))

		res, err := tool.Execute(tmpDir, []byte(`{"path":"."}`))

		require.NoError(t, err)
		require.Contains(t, res.Stdout, "a.txt")
		require.Contains(t, res.Stdout, "b.txt")
	})
}

func TestMkdir_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success: create dir", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewMkdir()

		res, err := tool.Execute(tmpDir, []byte(`{"path":"a/b/c"}`))

		require.NoError(t, err)
		require.NotNil(t, res)

		info, err := os.Stat(filepath.Join(tmpDir, "a/b/c"))
		require.NoError(t, err)
		require.True(t, info.IsDir())
	})

	t.Run("error: empty path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewMkdir()

		_, err := tool.Execute(tmpDir, []byte(`{"path":""}`))

		require.Error(t, err)
	})
}

func TestFileMove_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success: move file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewFileMove()

		old := filepath.Join(tmpDir, "a.txt")
		new := filepath.Join(tmpDir, "b.txt")

		require.NoError(t, os.WriteFile(old, []byte("x"), 0644))

		res, err := tool.Execute(tmpDir, []byte(`{"from":"a.txt","to":"b.txt"}`))

		require.NoError(t, err)
		require.NotNil(t, res)

		_, err = os.Stat(new)
		require.NoError(t, err)
	})

	t.Run("error: missing file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewFileMove()

		_, err := tool.Execute(tmpDir, []byte(`{"from":"a.txt","to":"b.txt"}`))

		require.Error(t, err)
	})
}
