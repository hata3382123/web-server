package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, authorId int64, status domain.ArticleStatus) error
}

type articleRepository struct {
	dao dao.ArticleDao
}

func NewArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &articleRepository{
		dao: dao,
	}
}
func (r *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return r.dao.Insert(ctx, dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}
func (r *articleRepository) Update(ctx context.Context, art domain.Article) error {
	return r.dao.Update(ctx, dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}

func (r *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return r.dao.SyncV2(ctx, dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}
func (r *articleRepository) SyncStatus(ctx context.Context, id int64, authorId int64, status domain.ArticleStatus) error {
	return r.dao.SyncStatus(ctx, id, authorId, status.ToUint8())
}
