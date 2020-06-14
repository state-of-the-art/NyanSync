package processors

type FileSystemType string

const (
	Source FileSystemType = "source"
	Folder FileSystemType = "folder"
	Binary FileSystemType = "binary"
)

type FileSystemItem struct {
	Name string
	Type FileSystemType
}
