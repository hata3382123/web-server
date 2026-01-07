package repository

import (
	"context"
	"database/sql"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, userId int64) (domain.User, error)
	Update(ctx context.Context, u domain.User) error
}
type CacheUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, c cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}
func (r *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}
func (r *CacheUserRepository) Update(ctx context.Context, u domain.User) error {
	return r.dao.Update(ctx, dao.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
	})
}
func (r *CacheUserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {
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
	u = r.entityToDomain(ue)
	err = r.cache.Set(ctx, u)
	return u, nil
}
func (r *CacheUserRepository) domainToEntity(u domain.User) dao.User {
	var ctime int64
	if !u.Ctime.IsZero() {
		ctime = u.Ctime.UnixMilli()
	}
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Ctime: ctime,
	}
}

func (r *CacheUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

//缓存里没数据
//缓存出错
