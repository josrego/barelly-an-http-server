package encodings

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
)

type Encoder struct {
	EncodingType string
	EncodingFun  func(*[]byte) ([]byte, error)
}

var AvailableEncoders = map[string]Encoder{
	"gzip": {EncodingType: "gzip", EncodingFun: encodeGzip},
}

func encodeGzip(body *[]byte) ([]byte, error) {
	// Create a buffer for the compressed data
	var buf bytes.Buffer

	gzWriter := gzip.NewWriter(&buf)

	// Write the original data into the GZIP writer
	_, err := gzWriter.Write(*body)
	if err != nil {
		errorMsg := fmt.Sprintf("Error while writing body in gzip writter: %s", err.Error())
		return nil, errors.New(errorMsg)
	}

	err = gzWriter.Close()
	if err != nil {
		errorMsg := fmt.Sprintf("Error while closing gzip writter: %s", err.Error())
		return nil, errors.New(errorMsg)
	}

	return buf.Bytes(), nil
}
