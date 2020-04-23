package gui

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/state-of-the-art/NyanSync/lib/generated"
)

const (
	theme_default = "default"
	theme_prefix = "/theme-assets/"
	theme_gui_prefix = "gui/"
)

func Init(p string, router *gin.Engine) {
	// Use local FS or the embedded one
	if p != "" {
		dir := http.Dir(p)
		fs = &dir
	} else {
		fs = generated.AssetFile()
	}

	router.NoRoute(doRoute)
}

func doRoute(c *gin.Context) {
	fmt.Printf("[DEBUG] Processing static: %s\n", c.Request.URL.Path)

	p := c.Request.URL.Path
	file := path.Base(p)

	// Serve our index file
	if file == "" || filepath.Ext(file) == "" {
		c.FileFromFS("gui/default/index.htm", fs)
		fmt.Printf("[DEBUG] Processing static 2: %s, %s\n", file, c.Writer.Status())
		return
	}

	// Serve special theme
	if strings.HasPrefix(p, theme_prefix) {
		c.FileFromFS(strings.Replace(p, theme_prefix, theme_gui_prefix, 1), fs)
		return
	}

	theme := theme_default // TODO: use selected theme

	// Check the preset theme
	if _, err := fs.Open(theme_gui_prefix + theme + p); err == nil {
		c.FileFromFS(theme_gui_prefix + theme + p, fs)
		return
	}

	// Serve default theme
	c.FileFromFS(theme_gui_prefix + theme_default + p, fs)
}

var fs http.FileSystem
