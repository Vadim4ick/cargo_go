package file

import (
	"context"
	"io"
)

type Meta struct {
	URL string
}

type Storage interface {
	Save(ctx context.Context, name string, r io.Reader) (Meta, error)
	Delete(ctx context.Context, url string) error
}
