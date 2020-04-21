// Generated package contains automatically generated files
package generated

//go:generate go get github.com/go-bindata/go-bindata/v3/...
//go:generate go-bindata -fs -prefix "../../" -pkg "generated" -o gui.assets.generated.go ../../gui/...

import (
	"fmt"
	"strings"
	"net/http"
)

const (
	theme_default = "default"
	theme_prefix = "/theme-assets/"
)


func Gui(path string) http.FileSystem {
	if path != "" {
		server := http.Dir(path)
		return &guiAssetsFileSystem{ &server }
	}
	return &guiAssetsFileSystem{ &assetOperator{} }
}


type guiAssetsFileSystem struct {
	parent http.FileSystem
}

func (fs guiAssetsFileSystem) Open(path string) (http.File, error) {
	fmt.Printf("[DEBUG] Processing static: %s\n", path)
	if path == "/" {
		return fs.parent.Open("gui")
	}

	theme := theme_default // TODO: use selected theme

	// If it starts with theme prefix, then use specified path
	if strings.HasPrefix(path, theme_prefix) {
		return fs.parent.Open(strings.Replace(path, theme_prefix, "gui/", 1))
	}

	f, err := fs.parent.Open("gui/" + theme + path)
	if err != nil {
		f, err = fs.parent.Open("gui/" + theme_default + path)
	}

	return f, err
}
