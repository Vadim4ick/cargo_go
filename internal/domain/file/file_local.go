package file

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
)

type Local struct{ Dir, BaseURL string }

func (l Local) Save(ctx context.Context, name string, r io.Reader) (Meta, error) {
	fname := uuid.NewString() + filepath.Ext(name)
	full := filepath.Join(l.Dir, fname)

	if dst, err := os.Create(full); err != nil {
		return Meta{}, err
	} else {
		defer dst.Close()
		if _, err := io.Copy(dst, r); err != nil {
			return Meta{}, err
		}
	}

	return Meta{URL: path.Join(l.BaseURL, fname)}, nil
}

func (l Local) Delete(_ context.Context, url string) error {
	return os.Remove(filepath.Join(l.Dir, filepath.Base(url)))
}
