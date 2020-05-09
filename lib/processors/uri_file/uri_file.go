package uri_file

import (
	"net/url"
	"os"

	"github.com/pkg/errors"

	"github.com/state-of-the-art/NyanSync/lib/processors"
)

func init() {
	processors.Uri = append(processors.Uri, processors.UriProcessor{
		Scheme:  "file",
		IsValid: isValid,
		GetList: getList,
	})
}

func isValid(uri *url.URL) (err error) {
	_, err = os.Open(uri.Path)
	return
}

func getList(uri *url.URL) (data []processors.FileSystemItem, err error) {
	// TODO: mark invalid paths contains "/.." "../"
	var fh *os.File
	fh, err = os.Open(uri.Path)
	if err != nil {
		return
	}
	defer fh.Close()

	var list []os.FileInfo
	list, err = fh.Readdir(-1)
	if err != nil {
		err = errors.Wrap(err, processors.ErrUriUnableToList)
		return
	}

	for _, file := range list {
		fsi := processors.FileSystemItem{
			Name: file.Name(),
		}
		if file.IsDir() {
			fsi.Type = processors.Folder
		} else {
			fsi.Type = processors.Binary
		}
		data = append(data, fsi)
	}

	return
}
