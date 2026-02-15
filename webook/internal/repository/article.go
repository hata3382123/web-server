package repository

import (
	"context"
	"errors"
	"webook/internal/domain"
	dao "webook/internal/repository/dao/article"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

type CacheArticleRepository struct {
	dao dao.ArticleDao

	//操作两个DAO
	readerDao dao.ReaderDao
	authorDao dao.AuthorDao
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
	err := r.dao.Update(ctx, r.domainToEntity(art))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrArticleNotFound
	}
	return err
}
func (r *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	//先保存到制作库，再保存到线上库
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = r.authorDao.Update(ctx, r.domainToEntity(art))
	} else {
		id, err = r.authorDao.Insert(ctx, r.domainToEntity(art))
	}
	if err != nil {
		return 0, err
	}
	err = r.readerDao.Upsert(ctx, r.domainToEntity(art))
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (r *CacheArticleRepository) domainToEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
