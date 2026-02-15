package repository

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository/dao"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}

type CacheArticleRepository struct {
	dao dao.ArticleDao
}

func NewArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &CacheArticleRepository{
		dao: dao,
	}
}
func (r *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return r.dao.Insert(ctx, dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}

// ErrArticleNotFound 文章不存在或无权操作
var ErrArticleNotFound = errors.New("article not found or permission denied")

func (r *CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	err := r.dao.Update(ctx, dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrArticleNotFound
	}
	return err
}
