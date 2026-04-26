package tools

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistry_Register(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		r := NewRegistry()
		err := r.Register(mockTool{code: "code"})
		require.NoError(t, err)
	})

	t.Run("duplicate error", func(t *testing.T) {
		t.Parallel()

		r := NewRegistry()
		err := r.Register(mockTool{code: "code"})
		require.NoError(t, err)

		err = r.Register(mockTool{code: "code"})
		require.ErrorIs(t, err, ErrAlreadyRegistered)
	})
}

func TestRegistry_Get(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		tool := mockTool{code: "code"}

		r := NewRegistry()
		err := r.Register(tool)
		require.NoError(t, err)

		result, err := r.Get(tool.code)
		require.NoError(t, err)
		require.Equal(t, tool, result)
	})

	t.Run("not found error", func(t *testing.T) {
		t.Parallel()

		r := NewRegistry()

		result, err := r.Get("code")
		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, result)
	})
}

func TestRunCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cmd  *exec.Cmd
		want Result
		err  error
	}{
		{
			name: "success stdout only",
			cmd:  exec.Command("sh", "-c", "echo ok"),
			want: Result{
				Stdout: "ok\n",
				StdErr: "",
			},
		},
		{
			name: "stderr but success exit",
			cmd:  exec.Command("sh", "-c", "echo warn 1>&2; echo ok"),
			want: Result{
				Stdout: "ok\n",
				StdErr: "warn\n",
			},
		},
		{
			name: "failure with stderr",
			cmd:  exec.Command("sh", "-c", "echo fail 1>&2; exit 1"),
			want: Result{
				Stdout: "",
				StdErr: "fail\n",
			},
			err: ErrRunCommand,
		},
		{
			name: "empty output success",
			cmd:  exec.Command("sh", "-c", "exit 0"),
			want: Result{
				Stdout: "",
				StdErr: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := runCommand(tt.cmd)

			require.Equal(t, tt.want.Stdout, res.Stdout)
			require.Equal(t, tt.want.StdErr, res.StdErr)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

type mockTool struct {
	code string
}

func (m mockTool) Code() string {
	return m.code
}

func (m mockTool) Spec() Spec {
	return Spec{}
}

func (m mockTool) Execute(string, []byte) (Result, error) {
	return Result{}, nil
}
