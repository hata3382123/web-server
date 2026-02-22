package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
}

type ArticleStatus uint8

func (a ArticleStatus) ToUint8() uint8 {
	return uint8(a)
}

const (
	//unknow为了避免0值
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusPublished
	ArticleStatusUnPublished
	ArticleStatusPrivate
)

type Author struct {
	Name string
	Id   int64
}
