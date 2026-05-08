package git

import (
	"testing"

	"github.com/cyradin/govnocode/internal/command"
	"github.com/stretchr/testify/require"
)

func TestPush_Execute(t *testing.T) {
	t.Parallel()

	t.Run("expected failure without remote", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		e := command.NewShellExecutor(tmpDir)

		initRepo(t, tmpDir, e)

		_, err := NewPush().Execute(t.Context(), e, []byte(`{"remote":"origin","branch":"main"}`))

		require.Error(t, err)
	})
}
