package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyradin/govnocode/tools/executor"
	"github.com/stretchr/testify/require"
)

func TestAdd_Execute(t *testing.T) {
	t.Parallel()

	t.Run("add file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		initRepo(t, tmpDir, e)

		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("x"), 0644))

		res, err := NewAdd().Execute(t.Context(), e, []byte(`{"path":"."}`))

		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("invalid path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := executor.NewShell(tmpDir)

		initRepo(t, tmpDir, e)

		_, err := NewAdd().Execute(t.Context(), e, []byte(`{"path":""}`))

		require.Error(t, err)
	})
}
