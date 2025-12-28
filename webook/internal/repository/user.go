package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}
func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	return r.dao.Update(ctx, dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
	})
}
func (r *UserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {
	//u, err := r.dao.FindById(ctx, userId)
	//if err != nil {
	//	return domain.User{}, err
	//}

	//缓存里有数据
	u, err := r.cache.Get(ctx, userId)
	if err == nil {
		//有数据
		return u, err
	}
	//没这个数据
	//if err == cache.ErrKeyNotExist {
	//
	//}
	ue, err := r.dao.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
		Nickname: ue.Nickname,
		Birthday: ue.Birthday,
		AboutMe:  ue.AboutMe,
	}
	err = r.cache.Set(ctx, u)
	return u, nil
}

//缓存里没数据
//缓存出错
