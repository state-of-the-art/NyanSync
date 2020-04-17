package location

import (
    "fmt"
    "os"
    "strings"
    "runtime"
    "path/filepath"

    "github.com/pkg/errors"      // errors wrapper
)

func ExpandTilde(path string) (string, error) {
    if path == "~" {
        return os.UserHomeDir()
    }

    path = filepath.FromSlash(path)
    if !strings.HasPrefix(path, fmt.Sprintf("~%c", os.PathSeparator)) {
        return path, nil
    }

    home, err := os.UserHomeDir()
    if err != nil {
        return "", errors.Wrap(err, "user home dir")
    }
    return filepath.Join(home, path[2:]), nil
}

func DefaultConfigDir(app_name string) (dir string, err error) {
    switch runtime.GOOS {
    case "windows":
        if p := os.Getenv("LocalAppData"); p != "" {
            return filepath.Join(p, app_name), nil
        }
        dir = filepath.Join(os.Getenv("AppData"), app_name)

    case "darwin":
        dir, err = ExpandTilde(filepath.Join("~/Library/Application Support", app_name))

    default:
        if xdgCfg := os.Getenv("XDG_CONFIG_HOME"); xdgCfg != "" {
            return filepath.Join(xdgCfg, app_name), nil
        }
        dir, err = ExpandTilde(filepath.Join("~/.config", app_name))
    }
    return
}
