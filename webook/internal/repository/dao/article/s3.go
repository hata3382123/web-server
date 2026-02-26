package article

import "github.com/aws/aws-sdk-go-v2/service/s3"

type S3DAO struct {
	oss *s3.Client

	//通过组合GORMArticleDao来简化操作
	//当然在实践中 是不太会有组合的机会
	//你操作制作库总是一样的
	//就是操作线上库的时候不一样
	GORMArticleDao
	bucket *string
}
