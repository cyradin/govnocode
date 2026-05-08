package command

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecutor_Run_Success(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ex := NewExecutor(dir)

	res, err := ex.Run(context.Background(), []string{
		"sh", "-c", "echo hello",
	})

	require.NoError(t, err)
	require.Equal(t, "hello\n", res.Stdout)
	require.Equal(t, "", res.Stderr)
}

func TestExecutor_Run_WorkingDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	filePath := filepath.Join(dir, "test.txt")
	err := os.WriteFile(filePath, []byte("ok"), 0644)
	require.NoError(t, err)

	ex := NewExecutor(dir)

	res, err := ex.Run(context.Background(), []string{
		"cat", "test.txt",
	})

	require.NoError(t, err)
	require.Equal(t, "ok", res.Stdout)
}

func TestExecutor_Run_Error(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ex := NewExecutor(dir)

	_, err := ex.Run(context.Background(), []string{
		"sh", "-c", "exit 1",
	})

	require.Error(t, err)
	require.ErrorIs(t, err, ErrRunCommand)
}
