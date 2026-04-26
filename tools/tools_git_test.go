//go:build integration

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

		initRepo(t, tmpDir)

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

		initRepo(t, tmpDir)

		res, err := tool.Execute(tmpDir, []byte(`{"branch":"--help"}`))

		require.Error(t, err)
		require.NotNil(t, res)
	})

	t.Run("error: empty branch", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGetCreateBranch()

		initRepo(t, tmpDir)

		_, err := tool.Execute(tmpDir, []byte(`{"branch":""}`))

		require.Error(t, err)
	})
}

func TestGitCheckout_Execute(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	t.Run("success: switch branch", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitCheckout()

		initRepo(t, tmpDir)

		_, err := runCommand(exec.Command("git", "branch", "feature/test"), tmpDir)
		require.NoError(t, err)

		res, err := tool.Execute(tmpDir, []byte(`{"branch":"feature/test"}`))

		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("error: checkout non-existent branch", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitCheckout()

		initRepo(t, tmpDir)

		_, err := tool.Execute(tmpDir, []byte(`{"branch":"does-not-exist"}`))

		require.Error(t, err)
	})
}

func TestGitStatus_Execute(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	t.Run("success: empty repo status", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitStatus()

		initRepo(t, tmpDir)

		res, err := tool.Execute(tmpDir, nil)

		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("success: status after file change", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitStatus()

		initRepo(t, tmpDir)

		_ = os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("hello"), 0644)

		res, err := tool.Execute(tmpDir, nil)

		require.NoError(t, err)
		require.Contains(t, res.Stdout, "test.txt")
	})
}

func TestGitAdd_Execute(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	t.Run("success: add file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitAdd()

		initRepo(t, tmpDir)

		file := filepath.Join(tmpDir, "a.txt")
		require.NoError(t, os.WriteFile(file, []byte("data"), 0644))

		res, err := tool.Execute(tmpDir, []byte(`{"path":"."}`))

		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("error: invalid path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitAdd()

		initRepo(t, tmpDir)

		_, err := tool.Execute(tmpDir, []byte(`{"path":""}`))

		require.Error(t, err)
	})
}

func TestGitCommit_Execute(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	t.Run("success: commit file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitCommit()

		initRepo(t, tmpDir)

		_ = os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("data"), 0644)

		_, err := runCommand(exec.Command("git", "add", "."), tmpDir)
		require.NoError(t, err)

		res, err := tool.Execute(tmpDir, []byte(`{"message":"test commit"}`))

		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("error: empty message", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitCommit()

		initRepo(t, tmpDir)

		_, err := tool.Execute(tmpDir, []byte(`{"message":""}`))

		require.Error(t, err)
	})
}

func TestGitPush_Execute(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	t.Run("error: push without remote (expected failure)", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitPush()

		initRepo(t, tmpDir)

		_, err := tool.Execute(tmpDir, []byte(`{"remote":"origin","branch":"main"}`))

		require.Error(t, err)
	})
}

func TestGitPull_Execute(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	t.Run("error: pull without remote", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitPull()

		initRepo(t, tmpDir)

		_, err := tool.Execute(tmpDir, nil)

		require.Error(t, err)
	})
}

func TestGitFetch_Execute(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	t.Run("success: fetch in empty repo", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitFetch()

		_, err := exec.Command("git", "init", tmpDir).CombinedOutput()
		require.NoError(t, err)

		res, err := tool.Execute(tmpDir, nil)

		require.NoError(t, err)
		require.NotNil(t, res)
	})
}

func TestGitBranchList_Execute(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	t.Run("success: list branches", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		tool := NewGitBranchList()

		initRepo(t, tmpDir)

		_, err := runCommand(exec.Command("git", "branch", "feature/test"), tmpDir)
		require.NoError(t, err)

		res, err := tool.Execute(tmpDir, nil)

		require.NoError(t, err)
		require.Contains(t, res.Stdout, "test")
	})
}

func initRepo(t *testing.T, dir string) {
	t.Helper()

	_, _ = exec.Command("git", "-C", dir, "init").CombinedOutput()
	_, _ = exec.Command("git", "-C", dir, "config", "--local", "commit.gpgsign", "false").CombinedOutput()
	_, _ = exec.Command("git", "-C", dir, "config", "--local", "tag.gpgsign", "false").CombinedOutput()
	_, _ = exec.Command("git", "-C", dir, "config", "user.email", "test@test.com").CombinedOutput()
	_, _ = exec.Command("git", "-C", dir, "config", "user.name", "test").CombinedOutput()
	_, err := exec.Command("git", "-C", dir, "commit", "--allow-empty", "-m", "init").CombinedOutput()
	require.NoError(t, err)
}
