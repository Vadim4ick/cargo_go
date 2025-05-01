package file

import (
	"context"
	"io"
)

type Meta struct {
	URL string
}

type Record struct {
	ID, OwnerID, OwnerTable, URL string
}

type Repository interface {
	Create(ctx context.Context, rec Record) error
	DeleteByIDs(ctx context.Context, ids []string) ([]Record, error)
	GetByOwner(ctx context.Context, ownerTable, ownerID string) ([]Record, error)
}

type Storage interface {
	Save(ctx context.Context, name string, r io.Reader) (Meta, error)
	Delete(ctx context.Context, url string) error
}
