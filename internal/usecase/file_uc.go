package usecase

import (
	"context"
	"mime/multipart"
	"test-project/internal/domain/file"

	"github.com/google/uuid"
)

type FileService struct {
	st   file.Storage
	repo file.Repository
}

func NewFileService(st file.Storage, repo file.Repository) *FileService {
	return &FileService{st: st, repo: repo}
}

func (s *FileService) UploadMany(
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

		rec := file.Record{
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

func (s *FileService) DeleteMany(ctx context.Context, ids []string) error {
	recs, err := s.repo.DeleteByIDs(ctx, ids)
	if err != nil {
		return err
	}

	for _, r := range recs {
		_ = s.st.Delete(ctx, r.URL)
	}
	return nil
}
