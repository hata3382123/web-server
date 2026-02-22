package article

import (
	"context"

	"gorm.io/gorm"
)

type AuthorDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
}

func NewAuthorDao(db *gorm.DB) AuthorDao {
	return &GORMArticleDao{
		db: db,
	}
}
