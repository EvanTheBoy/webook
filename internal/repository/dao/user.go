package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
	"webook/internal/domain"
)

// ErrUserDuplicateEmail 每层都有自己的error, 每一层只能调它下一层的error, 不能越级调
var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

// UserDAO 数据库, 存储意义上的用户
type UserDAO struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// 针对存储意义上的User进行存储, 在这里主要存创建时间
	// 与更新时间, 然后调用context让它能一直在链路上传递下去
	// 这里我们只需要存时间的毫秒数, 用来抵消时区的影响, 只用
	// 在前端展示的时候再做一下处理就好了
	now := time.Now().UnixMilli()
	u.CreatedTime = now
	u.UpdatedTime = now

	err := dao.db.WithContext(ctx).Create(&u).Error
	// 邮箱冲突的错误具有唯一索引, 就是它的错误码是唯一的
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) SelectEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *UserDAO) UpdateById(ctx *gin.Context, u User) error {
	now := time.Now().UnixMilli()
	result := dao.db.WithContext(ctx).Model(&u).Where("id = ?", u.Id).Updates(map[string]interface{}{
		"Nickname":    u.Nickname,
		"Birthday":    u.Birthday,
		"Address":     u.Address,
		"BriefIntro":  u.BriefIntro,
		"UpdatedTime": now,
	})
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return result.Error
}

func (dao *UserDAO) SelectUserById(ctx *gin.Context, u domain.User) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id = ?", u.Id).First(&user).Error
	return user, err
}

func (dao *UserDAO) SelectPhone(ctx *gin.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

type User struct {
	Id         int64          `gorm:"primaryKey,autoIncrement"`
	Email      sql.NullString `gorm:"unique"`
	Phone      sql.NullString `gorm:"unique"`
	Password   string
	Nickname   string `gorm:"type=varchar(20)"`
	Birthday   string
	Address    string `gorm:"type=varchar(40)"`
	BriefIntro string `gorm:"type=varchar(60)"`

	CreatedTime int64
	UpdatedTime int64
}
