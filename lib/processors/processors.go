package processors

type FileSystemType string

const (
	Folder FileSystemType = "folder"
	Binary FileSystemType = "binary"
)

type FileSystemItem struct {
	Name string
	Type FileSystemType
}
