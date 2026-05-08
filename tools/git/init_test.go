package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cyradin/govnocode/tools/executor"
	"github.com/stretchr/testify/require"
)

func TestInit_Execute(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	t.Run("success: initializes git repository", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		executor := executor.NewShell(tmpDir)

		tool := NewInit()

		res, err := tool.Execute(t.Context(), executor, nil)

		require.NoError(t, err)
		require.NotNil(t, res)

		gitDir := filepath.Join(tmpDir, ".git")
		info, err := os.Stat(gitDir)
		require.NoError(t, err)
		require.True(t, info.IsDir())
	})

	t.Run("success: idempotent (second init does not fail)", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		executor := executor.NewShell(tmpDir)

		tool := NewInit()

		_, err := tool.Execute(t.Context(), executor, nil)
		require.NoError(t, err)

		_, err = tool.Execute(t.Context(), executor, nil)
		require.NoError(t, err)
	})

	t.Run("error: invalid directory", func(t *testing.T) {
		t.Parallel()

		tool := NewInit()

		executor := executor.NewShell("/non/existent/dir")

		_, err := tool.Execute(t.Context(), executor, nil)
		require.Error(t, err)
	})
}

func initRepo(t *testing.T, dir string, executor executor.Executor) {
	t.Helper()

	_, _ = executor.Execute(t.Context(), []string{"git", "-C", dir, "init"}, nil)
	_, _ = executor.Execute(t.Context(), []string{"git", "-C", dir, "config", "--local", "user.email", "test@test.com"}, nil)
	_, _ = executor.Execute(t.Context(), []string{"git", "-C", dir, "config", "--local", "user.name", "test"}, nil)
	_, _ = executor.Execute(t.Context(), []string{"git", "-C", dir, "config", "--local", "commit.gpgsign", "false"}, nil)
	_, _ = executor.Execute(t.Context(), []string{"git", "-C", dir, "config", "--local", "tag.gpgsign", "false"}, nil)
	_, _ = executor.Execute(t.Context(), []string{"git", "-C", dir, "commit", "--allow-empty", "-m", "init"}, nil)
}
