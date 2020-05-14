package location

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	AppName = "NyanShare"
)

func ExpandTilde(path string) string {
	path = filepath.FromSlash(path)
	if !strings.HasPrefix(path, fmt.Sprintf("~%c", os.PathSeparator)) {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic("Unable to get user home dir", err)
	}
	return filepath.Join(home, path[2:])
}

func DefaultConfigDir() (dir string) {
	switch runtime.GOOS {
	case "windows":
		if p := os.Getenv("LocalAppData"); p != "" {
			return filepath.Join(p, AppName)
		}
		dir = filepath.Join(os.Getenv("AppData"), AppName)

	case "darwin":
		dir = ExpandTilde(filepath.Join("~/Library/Application Support", AppName))

	default:
		if xdg_cfg := os.Getenv("XDG_CONFIG_HOME"); xdg_cfg != "" {
			return filepath.Join(xdg_cfg, AppName)
		}
		dir = ExpandTilde(filepath.Join("~/.config", AppName))
	}
	return
}

func RealFilePath(file_path string) (path string) {
	path = ExpandTilde(file_path)
	if filepath.IsAbs(path) {
		return
	}

	default_dir := DefaultConfigDir()

	return filepath.Join(default_dir, path)
}
