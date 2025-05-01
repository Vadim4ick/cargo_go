package file

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
)

type Service struct {
	st   Storage
	repo Repo
}

func NewService(st Storage, repo Repo) *Service {
	return &Service{st: st, repo: repo}
}

func (s *Service) UploadMany(
	ctx context.Context,
	ownerTable, ownerID string,
	fhs []*multipart.FileHeader,
) error {
	for _, fh := range fhs {
		src, _ := fh.Open()
		defer src.Close()

		meta, err := s.st.Save(ctx, fh.Filename, src)
		if err != nil {
			return err
		}

		rec := Record{
			ID:         uuid.NewString(),
			OwnerID:    ownerID,
			OwnerTable: ownerTable,
			URL:        meta.URL,
		}
		if err := s.repo.Create(ctx, rec); err != nil {
			_ = s.st.Delete(ctx, meta.URL)
			return err
		}
	}
	return nil
}

func (s *Service) DeleteMany(ctx context.Context, ids []string) error {
	recs, err := s.repo.DeleteByIDs(ctx, ids)
	if err != nil {
		return err
	}

	for _, r := range recs {
		_ = s.st.Delete(ctx, r.URL)
	}
	return nil
}
