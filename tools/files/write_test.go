package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyradin/govnocode/tools/executor"
	"github.com/stretchr/testify/require"
)

func TestWrite_Execute(t *testing.T) {
	t.Parallel()

	t.Run("write new file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		res, err := NewWrite().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"a.txt","content":"hello"}`),
		)

		require.NoError(t, err)
		require.NotNil(t, res)

		b, err := os.ReadFile(filepath.Join(tmpDir, "a.txt"))
		require.NoError(t, err)

		require.Equal(t, "hello", string(b))
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		path := filepath.Join(tmpDir, "a.txt")

		require.NoError(t, os.WriteFile(path, []byte("old"), 0644))

		_, err := NewWrite().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"a.txt","content":"new content"}`),
		)

		require.NoError(t, err)

		b, err := os.ReadFile(path)
		require.NoError(t, err)

		require.Equal(t, "new content", string(b))
	})

	t.Run("empty content", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		_, err := NewWrite().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"a.txt","content":""}`),
		)

		require.NoError(t, err)

		b, err := os.ReadFile(filepath.Join(tmpDir, "a.txt"))
		require.NoError(t, err)

		require.Empty(t, string(b))
	})

	t.Run("missing path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		_, err := NewWrite().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"","content":"hello"}`),
		)

		require.Error(t, err)
	})

	t.Run("nested directory auto create", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		res, err := NewWrite().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"nested/a.txt","content":"hello"}`),
		)

		require.NoError(t, err)
		require.NotNil(t, res)

		b, err := os.ReadFile(filepath.Join(tmpDir, "nested/a.txt"))
		require.NoError(t, err)

		require.Equal(t, "hello", string(b))
	})

	t.Run("write multiline content", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		content := "line1\nline2\nline3"

		_, err := NewWrite().Execute(
			t.Context(),
			e,
			[]byte(`{"path":"a.txt","content":"line1\nline2\nline3"}`),
		)

		require.NoError(t, err)

		b, err := os.ReadFile(filepath.Join(tmpDir, "a.txt"))
		require.NoError(t, err)

		require.Equal(t, content, string(b))
	})
}
