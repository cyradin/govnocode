package executor

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShell_Execute_Success(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ex := NewShell(dir)

	res, err := ex.Execute(context.Background(), []string{
		"sh", "-c", "echo hello",
	}, nil)

	require.NoError(t, err)
	require.Equal(t, "hello\n", res.Stdout)
	require.Equal(t, "", res.Stderr)
}

func TestShell_Execute_WorkingDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	filePath := filepath.Join(dir, "test.txt")
	err := os.WriteFile(filePath, []byte("ok"), 0644)
	require.NoError(t, err)

	ex := NewShell(dir)

	res, err := ex.Execute(context.Background(), []string{
		"cat", "test.txt",
	}, nil)

	require.NoError(t, err)
	require.Equal(t, "ok", res.Stdout)
}

func TestShell_Execute_Error(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ex := NewShell(dir)

	_, err := ex.Execute(context.Background(), []string{
		"sh", "-c", "exit 1",
	}, nil)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrRunCommand)
}
