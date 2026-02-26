package article

import (
	"context"
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBDAO struct {
	// client *mongo.Client
	// //代表 webook的
	// database *mongo.Database
	//代表制作库
	col *mongo.Collection
	//代表线上库
	liveCol *mongo.Collection

	node *snowflake.Node
}

func (m *MongoDBDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	id := m.node.Generate().Int64()
	art.Id = id
	_, err := m.col.InsertOne(ctx, art)

	//没有自增主键
	return art.Id, err
}
func NewMongoDAO(db *mongo.Database, node *snowflake.Node) ArticleDao {
	return &MongoDBDAO{
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		node:    node,
	}
}

// SyncStatus implements ArticleDao.
func (*MongoDBDAO) SyncStatus(ctx context.Context, id int64, authorId int64, status uint8) error {
	panic("unimplemented")
}

// SyncV2 implements ArticleDao.
func (m *MongoDBDAO) SyncV2(ctx context.Context, art Article) (int64, error) {
	//操作制作库
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.Update(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}

	// 操作线上库：按 id upsert 到 published_articles 集合
	pub := PublishArticle{
		Article: art,
	}
	liveFilter := bson.M{"id": id}
	liveUpdate := bson.M{"$set": pub}
	_, err = m.liveCol.UpdateOne(ctx, liveFilter, liveUpdate, options.Update().SetUpsert(true))
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update implements ArticleDao.
func (m *MongoDBDAO) Update(ctx context.Context, art Article) error {
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	update := bson.D{bson.E{
		Key: "$set",
		Value: bson.M{
			"title":   art.Title,
			"content": art.Content,
			"utime":   time.Now().UnixMilli(),
			"status":  art.Status,
		},
	}}
	res, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

// Upsert implements ArticleDao.
func (*MongoDBDAO) Upsert(ctx context.Context, art PublishArticle) error {
	panic("unimplemented")
}
