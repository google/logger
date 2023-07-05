//go:build !unix && !windows

package logger

import (
	"errors"
	"io"
)

func setup(src string) (io.Writer, io.Writer, io.Writer, error) {
	return nil, nil, nil, errors.New("system logging not implemented")
}
