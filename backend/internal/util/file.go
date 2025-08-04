package util

import (
	"bytes"
	"io"
	"net/http"
)

// DetectContentType uses http.DetectContentType to determine the mime-type of a given io.Reader
func DetectContentType(r io.Reader) (string, io.Reader, error) {
	buf := make([]byte, 512)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", nil, err
	}
	buf = buf[:n]
	contentType := http.DetectContentType(buf)
	return contentType, io.MultiReader(bytes.NewReader(buf), r), nil
}
