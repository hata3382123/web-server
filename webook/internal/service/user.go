package service

import (
	"context"
	"errors"
	"webook/internal/domain"
	"webook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("账号或密码不正确")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.User, error)
	Edit(ctx context.Context, u domain.User) error
	Profile(ctx context.Context, userId int64) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	// 考虑加密放哪里
	// 存起来
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}
func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	//判断是否有用户
	if err == nil {
		// 找到了用户
		return u, nil
	}
	if err != repository.ErrUserNotFound {
		// 其他错误
		return domain.User{}, err
	}
	// 用户不存在，创建新用户
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != ErrUserDuplicate {
		return u, err
	}
	//这里会遇到主从延迟的问题，如果查询不到，重试一次
	u, err = svc.repo.FindByPhone(ctx, phone)
	if err != nil {
		// 如果还是找不到，可能是主从延迟，返回创建的用户（但ID可能为0）
		// 或者返回错误
		return domain.User{}, err
	}
	return u, nil
}
func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}
func (svc *userService) Edit(ctx context.Context, u domain.User) error {
	return svc.repo.Update(ctx, u)
}
func (svc *userService) Profile(ctx context.Context, userId int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	return u, err
}
