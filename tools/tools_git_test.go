//go:build integration
// +build integration

package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitInit_Execute(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	t.Run("success: initializes git repository", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitInit()

		res, err := tool.Execute(tmpDir, nil)

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
		tool := NewGitInit()

		_, err := tool.Execute(tmpDir, nil)
		require.NoError(t, err)

		_, err = tool.Execute(tmpDir, nil)
		require.NoError(t, err)
	})

	t.Run("error: invalid directory", func(t *testing.T) {
		t.Parallel()

		tool := NewGitInit()

		_, err := tool.Execute("/non/existent/dir", nil)
		require.Error(t, err)
	})
}

func TestGitCreateBranch_Execute(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	t.Run("success: create new branch", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGetCreateBranch()

		_, err := exec.Command("git", "init", tmpDir).CombinedOutput()
		require.NoError(t, err)

		res, err := tool.Execute(tmpDir, []byte(`{"branch":"feature/test"}`))

		require.NoError(t, err)
		require.NotNil(t, res)

		out, err := exec.Command("git", "-C", tmpDir, "branch", "--show-current").Output()
		require.NoError(t, err)

		require.Equal(t, "feature/test\n", string(out))
	})

	t.Run("error: invalid branch name", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGetCreateBranch()

		_, err := exec.Command("git", "init", tmpDir).CombinedOutput()
		require.NoError(t, err)

		res, err := tool.Execute(tmpDir, []byte(`{"branch":"--help"}`))

		require.Error(t, err)
		require.Contains(t, err.Error(), "git")
		require.NotNil(t, res)
	})

	t.Run("error: empty branch", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGetCreateBranch()

		_, err := exec.Command("git", "init", tmpDir).CombinedOutput()
		require.NoError(t, err)

		_, err = tool.Execute(tmpDir, []byte(`{"branch":""}`))

		require.Error(t, err)
	})
}
