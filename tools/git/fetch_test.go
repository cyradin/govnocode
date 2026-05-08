package git

import (
	"testing"

	"github.com/cyradin/govnocode/internal/command"
	"github.com/stretchr/testify/require"
)

func TestFetch_Execute(t *testing.T) {
	t.Parallel()

	t.Run("fetch empty repo", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		_, err := e.Execute(t.Context(), []string{"git", "-C", tmpDir, "init"}, nil)
		require.NoError(t, err)

		res, err := NewFetch().Execute(t.Context(), e, nil)

		require.NoError(t, err)
		require.NotNil(t, res)
	})
}
