package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicate = errors.New("邮箱或手机号冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}
func (dao *UserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	//存毫秒数
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			//邮箱或手机号码冲突
			return ErrUserDuplicate
		}
	}
	return err
}
func (dao *UserDao) Update(ctx context.Context, u User) error {
	//存毫秒数
	now := time.Now().UnixMilli()
	// 只更新指定字段，避免零值覆盖
	updates := map[string]interface{}{
		"utime": now,
	}
	if u.Nickname != "" {
		updates["nickname"] = u.Nickname
	}
	if u.Birthday != "" {
		updates["birthday"] = u.Birthday
	}
	if u.AboutMe != "" {
		updates["about_me"] = u.AboutMe
	}
	err := dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", u.Id).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}
func (dao *UserDao) FindById(ctx context.Context, userId int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", userId).First(&u).Error
	return u, err
}

// user 直接对应 数据库
// 有些人叫做entity 有些人叫model 有些人叫做PO
type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`
	//往这里加需要的字段

	//创建时间 毫秒数
	Ctime int64

	//更新时间 毫秒数
	Utime    int64
	Nickname string `gorm:"column:nickname"`
	Birthday string `gorm:"column:birthday"`
	AboutMe  string `gorm:"column:about_me"`
}
