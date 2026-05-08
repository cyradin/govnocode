package git

import (
	"testing"

	"github.com/cyradin/govnocode/internal/command"
	"github.com/stretchr/testify/require"
)

func TestPull_Execute(t *testing.T) {
	t.Parallel()

	t.Run("pull without remote", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		initRepo(t, tmpDir, e)

		_, err := NewPull().Execute(t.Context(), e, nil)

		require.Error(t, err)
	})
}
