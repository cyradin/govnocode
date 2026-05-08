package tools

import (
	"github.com/cyradin/govnocode/tools/files"
)

func Files() []*Tool {
	return []*Tool{
		NewFileRead(),
		NewFileWrite(),
		NewFileDelete(),
		NewDirList(),
		NewMkdir(),
		NewFileMove(),
	}
}

func NewFileRead() *Tool {
	return NewTool(
		files.NewRead(),
		Spec{
			Code:        "file.read",
			Description: "Read file content",
			Args: files.ReadArgs{
				Path: "path/to/file.txt",
			},
		},
	)
}

func NewFileWrite() *Tool {
	return NewTool(
		files.NewWrite(),
		Spec{
			Code:        "file.write",
			Description: "Write content to file (overwrite if exists)",
			Args: files.WriteArgs{
				Path:    "path/to/file.txt",
				Content: "hello world",
			},
		},
	)
}

func NewFileDelete() *Tool {
	return NewTool(
		files.NewDelete(),
		Spec{
			Code:        "file.delete",
			Description: "Delete file or directory",
			Args: files.DeleteArgs{
				Path: "path/to/file.txt",
			},
		},
	)
}

func NewDirList() *Tool {
	return NewTool(
		files.NewDirList(),
		Spec{
			Code:        "dir.list",
			Description: "List current directory contents",
			Args: files.DirListArgs{
				Path: ".",
			},
		},
	)
}

func NewMkdir() *Tool {
	return NewTool(
		files.NewMkdir(),
		Spec{
			Code:        "dir.mkdir",
			Description: "Create directory (recursive)",
			Args: files.MkdirArgs{
				Path: "path/to/dir",
			},
		},
	)
}

func NewFileMove() *Tool {
	return NewTool(
		files.NewMove(),
		Spec{
			Code:        "file.move",
			Description: "Move or rename file/directory",
			Args: files.MoveArgs{
				From: "old/path.txt",
				To:   "new/path.txt",
			},
		},
	)
}
