package multipart

import (
	"bytes"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewFormFile(mw *multipart.Writer, key, filename string, data []byte) error {
	w, err := mw.CreateFormFile(key, filename)
	if err != nil {
		return err
	}

	if _, err = w.Write(data); err != nil {
		return err
	}
	if err = mw.Close(); err != nil {
		return err
	}

	return nil
}

func NewFormData() (*multipart.Form, error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	if err := NewFormFile(mw, "foo", "test1.txt", []byte("hello")); err != nil {
		return nil, err
	}
	if err := NewFormFile(mw, "foo", "test2.txt", []byte("world")); err != nil {
		return nil, err
	}
	if err := NewFormFile(mw, "bar", "test3.txt", []byte("hi")); err != nil {
		return nil, err
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}

	r := multipart.NewReader(&buf, mw.Boundary())
	return r.ReadForm(1 << 20)
}

func TestGetUploadsFromMultipart(t *testing.T) {
	form, err := NewFormData()
	require.NoError(t, err)

	for _, files := range form.File {
		mp := NewMultipartFromUploads(files)
		if assert.NotNil(t, mp) {
			assert.Equal(t, files, GetUploadsFromMultipart(mp))
		}
	}
}
