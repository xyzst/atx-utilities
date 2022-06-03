package util

import (
	"compress/gzip"
	"io"
	"net/http"
)

func DecompressResponse(response http.Response) (io.ReadCloser, error) {
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		r, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		reader = r
	default:
		reader = response.Body
	}

	return reader, nil
}
