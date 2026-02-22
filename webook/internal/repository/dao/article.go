package dao

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
	SyncV2(ctx context.Context, art Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, authorId int64, status uint8) error
}

type GORMArticleDao struct {
	db *gorm.DB
}

func NewGORMArticleDao(db *gorm.DB) ArticleDao {
	return &GORMArticleDao{
		db: db,
	}
}

func (dao *GORMArticleDao) Insert(ctx context.Context, art Article) (int64, error) {
	//return dao.db.WithContext(ctx).Create(&art).RowsAffected, nil
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (dao *GORMArticleDao) Update(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	// 只更新 title、content、utime，且必须 id + author_id 匹配，防止越权
	res := dao.db.WithContext(ctx).Model(&Article{}).
		Where("id = ? AND author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"utime":   now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		zap.L().Error("非法用户", zap.Error(res.Error))
		return gorm.ErrRecordNotFound // 文章不存在或不属于当前用户
	}
	return nil
}

func (dao *GORMArticleDao) SyncV2(ctx context.Context, art Article) (int64, error) {
	var id = art.Id
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMArticleDao(tx)
		if id > 0 {
			err := txDAO.Update(ctx, art)
			if err != nil {
				return err
			}
		} else {
			var err error
			id, err = txDAO.Insert(ctx, art)
			if err != nil {
				return err
			}
		}
		// 操作线上库（这里简化处理，实际可能需要 Upsert 到另一张表）
		return nil
	})
	return id, err
}

func (dao *GORMArticleDao) SyncStatus(ctx context.Context, id int64, authorId int64, status uint8) error {
	res := dao.db.WithContext(ctx).Model(&Article{}).
		Where("id = ? AND author_id = ?", id, authorId).
		Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	//长度1024
	Title   string `gorm:"type:varchar(1024)"`
	Content string `gorm:"type:BLOB"`
	//如果设计索引
	//在帖子这里，什么样的查询场景？
	//对于创作者来说，是不是草稿箱，看到所有自己的文章？
	//产品经理告诉你，要按照创建时间的倒序排序
	//SELECT * FROM articles WHERE author_id = 123 ORDER BY `ctime` DESC
	//单独查询某一篇 SELECT * FROM articles WHERE id = 1
	//在查询接口深入讨论这个问题
	// - 在 author_id 和 ctime 上建立联合索引
	// - 在author_id上建立索引

	//在author_id上建立索引
	AuthorId int64 `gorm:"index"`
	//在ctime上建立索引
	Status uint8
	Ctime  int64
	Utime  int64
}
