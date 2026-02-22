package article

import (
	"context"

	"gorm.io/gorm"
)

type ReaderDao interface {
	Upsert(ctx context.Context, art PublishArticle) error
}

func NewReaderDao(db *gorm.DB) ReaderDao {
	return &GORMArticleDao{
		db: db,
	}
}

type PublishArticle struct {
	Article
}
