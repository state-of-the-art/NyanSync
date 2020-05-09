package processors

import (
	"net/url"

	"github.com/pkg/errors"
)

const (
	ErrUriProcessorNotFound = "URI processor not found"
	ErrUriInvalid           = "URI is invalid"
	ErrUriUnableToList      = "Unable to list URI"
)

type UriProcessor struct {
	Scheme  string
	IsValid func(uri *url.URL) error
	GetList func(uri *url.URL) ([]FileSystemItem, error)
}

var Uri []UriProcessor

func findUriProcessor(uri *url.URL) (processor UriProcessor, err error) {
	for _, p := range Uri {
		if p.Scheme == uri.Scheme {
			processor = p
			return
		}
	}
	err = errors.New(ErrUriProcessorNotFound)
	return
}

func UriIsValid(uri *url.URL) error {
	proc, err := findUriProcessor(uri)
	if err != nil {
		return errors.Cause(err)
	}
	return proc.IsValid(uri)
}

func UriGetList(uri *url.URL) ([]FileSystemItem, error) {
	proc, err := findUriProcessor(uri)
	if err != nil {
		return nil, errors.Cause(err)
	}
	return proc.GetList(uri)
}
