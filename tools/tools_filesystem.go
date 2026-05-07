package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Filesystem() []Tool {
	return []Tool{
		NewFileRead(),
		NewFileWrite(),
		NewFileDelete(),
		NewDirList(),
		NewMkdir(),
		NewFileMove(),
	}
}

type fileReadArgs struct {
	Path string `json:"path" validate:"required"`
}

type FileRead struct{}

func NewFileRead() *FileRead {
	return &FileRead{}
}

func (f *FileRead) Code() string {
	return "file.read"
}

func (f *FileRead) Spec() Spec {
	return Spec{
		Code:        f.Code(),
		Description: "Read file content",
		Args: fileReadArgs{
			Path: "path/to/file.txt",
		},
	}
}

func (f *FileRead) Execute(dir string, raw []byte) (Result, error) {
	a, err := parseArgs[fileReadArgs](raw)
	if err != nil {
		return Result{}, fmt.Errorf("parse args: %w", err)
	}

	joinedPath, err := safeJoin(dir, a.Path)
	if err != nil {
		return Result{}, fmt.Errorf("path join: %w", err)
	}

	// #nosec G304
	b, err := os.ReadFile(joinedPath)
	if err != nil {
		return Result{}, fmt.Errorf("read file: %w", err)
	}

	return Result{Stdout: string(b)}, nil
}

type fileWriteArgs struct {
	Path    string `json:"path" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type FileWrite struct{}

func NewFileWrite() *FileWrite {
	return &FileWrite{}
}

func (f *FileWrite) Code() string {
	return "file.write"
}

func (f *FileWrite) Spec() Spec {
	return Spec{
		Code:        f.Code(),
		Description: "Write content to file (overwrite if exists)",
		Args: fileWriteArgs{
			Path:    "path/to/file.txt",
			Content: "hello world",
		},
	}
}

func (f *FileWrite) Execute(dir string, raw []byte) (Result, error) {
	a, err := parseArgs[fileWriteArgs](raw)
	if err != nil {
		return Result{}, fmt.Errorf("parse args: %w", err)
	}

	joinedPath, err := safeJoin(dir, a.Path)
	if err != nil {
		return Result{}, fmt.Errorf("path join: %w", err)
	}

	err = os.WriteFile(joinedPath, []byte(a.Content), 0600) //nolint:mnd
	if err != nil {
		return Result{}, fmt.Errorf("write file: %w", err)
	}

	return Result{Stdout: "ok"}, nil
}

type fileDeleteArgs struct {
	Path string `json:"path" validate:"required"`
}

type FileDelete struct{}

func NewFileDelete() *FileDelete {
	return &FileDelete{}
}

func (f *FileDelete) Code() string {
	return "file.delete"
}

func (f *FileDelete) Spec() Spec {
	return Spec{
		Code:        f.Code(),
		Description: "Delete file",
		Args: fileDeleteArgs{
			Path: "path/to/file.txt",
		},
	}
}

func (f *FileDelete) Execute(dir string, raw []byte) (Result, error) {
	a, err := parseArgs[fileDeleteArgs](raw)
	if err != nil {
		return Result{}, fmt.Errorf("parse args: %w", err)
	}

	joinedPath, err := safeJoin(dir, a.Path)
	if err != nil {
		return Result{}, fmt.Errorf("path join: %w", err)
	}

	err = os.Remove(joinedPath)
	if err != nil {
		return Result{}, fmt.Errorf("remove file: %w", err)
	}

	return Result{Stdout: "deleted"}, nil
}

type dirListArgs struct {
	Path string `json:"path"`
}

type DirList struct{}

func NewDirList() *DirList {
	return &DirList{}
}

func (d *DirList) Code() string {
	return "dir.list"
}

func (d *DirList) Spec() Spec {
	return Spec{
		Code:        d.Code(),
		Description: "List directory contents",
		Args: dirListArgs{
			Path: ".",
		},
	}
}

func (d *DirList) Execute(dir string, raw []byte) (Result, error) {
	a, err := parseArgs[dirListArgs](raw)
	if err != nil {
		return Result{}, fmt.Errorf("parse args: %w", err)
	}

	joinedPath, err := safeJoin(dir, a.Path)
	if err != nil {
		return Result{}, fmt.Errorf("path join: %w", err)
	}

	entries, err := os.ReadDir(joinedPath)
	if err != nil {
		return Result{}, fmt.Errorf("read dir: %w", err)
	}

	var b strings.Builder

	for _, e := range entries {
		b.WriteString(e.Name())
		b.WriteByte('\n')
	}

	return Result{Stdout: b.String()}, nil
}

type MkdirArgs struct {
	Path string `json:"path" validate:"required"`
}

type Mkdir struct{}

func NewMkdir() *Mkdir {
	return &Mkdir{}
}

func (m *Mkdir) Code() string {
	return "dir.mkdir"
}

func (m *Mkdir) Spec() Spec {
	return Spec{
		Code:        m.Code(),
		Description: "Create directory (recursive)",
		Args: MkdirArgs{
			Path: "path/to/dir",
		},
	}
}

func (m *Mkdir) Execute(dir string, raw []byte) (Result, error) {
	a, err := parseArgs[MkdirArgs](raw)
	if err != nil {
		return Result{}, fmt.Errorf("parse args: %w", err)
	}

	joinedPath, err := safeJoin(dir, a.Path)
	if err != nil {
		return Result{}, fmt.Errorf("path join: %w", err)
	}

	err = os.MkdirAll(joinedPath, 0750) //nolint:mnd
	if err != nil {
		return Result{}, fmt.Errorf("mkdir: %w", err)
	}

	return Result{Stdout: "ok"}, nil
}

type MoveArgs struct {
	From string `json:"from" validate:"required"`
	To   string `json:"to" validate:"required"`
}

type FileMove struct{}

func NewFileMove() *FileMove {
	return &FileMove{}
}

func (m *FileMove) Code() string {
	return "file.move"
}

func (m *FileMove) Spec() Spec {
	return Spec{
		Code:        m.Code(),
		Description: "Move or rename file/directory",
		Args: MoveArgs{
			From: "old/path.txt",
			To:   "new/path.txt",
		},
	}
}

func (m *FileMove) Execute(dir string, raw []byte) (Result, error) {
	a, err := parseArgs[MoveArgs](raw)
	if err != nil {
		return Result{}, fmt.Errorf("parse args: %w", err)
	}

	err = os.Rename(
		filepath.Join(dir, a.From),
		filepath.Join(dir, a.To),
	)
	if err != nil {
		return Result{}, fmt.Errorf("rename: %w", err)
	}

	return Result{Stdout: "moved"}, nil
}

func safeJoin(base, userPath string) (string, error) {
	full := filepath.Join(base, userPath)

	base = filepath.Clean(base)
	full = filepath.Clean(full)

	if full != base && !strings.HasPrefix(full, base+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid path: escapes base directory")
	}

	return full, nil
}
