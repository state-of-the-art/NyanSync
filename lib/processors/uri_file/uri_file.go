package uri_file

import (
	"io"
	"log"
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
		GetFile: getFile,
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

func getFile(uri *url.URL) (file processors.FileData, err error) {
	var fh *os.File
	fh, err = os.Open(uri.Path)
	if err != nil {
		return
	}

	var fi os.FileInfo
	fi, err = fh.Stat()
	if err != nil {
		err = errors.Wrap(err, processors.ErrUriUnableToStat)
		fh.Close()
		return
	}
	if fi.IsDir() {
		err = errors.Wrap(err, processors.ErrUriNotAFile)
		fh.Close()
		return
	}

	file = processors.FileData{
		Name: fi.Name(),
		Size: fi.Size(),
		Stream: func(ch chan []byte) {
			defer fh.Close()
			defer func() {
				// When channel is closed on read side
				if r := recover(); r != nil {
					log.Println("[WARN] Can't complete reading the file:", r)
				}
			}()
			data := make([]byte, 8192) // 8Kb buffer
			for {
				count, err := fh.Read(data)
				if count > 0 {
					ch <- data[:count]
				}
				if err != nil {
					if err == io.EOF {
						log.Println("[DEBUG] Reading the file compleated")
					} else {
						log.Println("[ERROR] Reading the file failed:", err)
					}
					close(ch)
					break
				}
			}
		},
	}

	return
}
