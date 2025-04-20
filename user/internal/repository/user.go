package repository

import (
	"errors"
	"etcd/user/internal/service"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	UserID          int    `gorm:"primarykey"`
	UserName        string `gorm:"unique"`
	NickName        string
	PasswordConfirm string
}

const (
	PasswordCost = 12 // 密码加密难度
)

// CheckUserExist 检查用户是否存在
func (user *User) CheckUserExist(req *service.UserRequest) bool {
	if err := DB.Where("user_name=?", req.UserName).First(&user).Error; err == gorm.ErrRecordNotFound {
		return false
	}

	return true
}

// ShowUserInfo 获取用户信息
func (user *User) ShowUserInfo(req *service.UserRequest) error {
	if exist := user.CheckUserExist(req); exist {
		return nil
	}

	return errors.New("UserName Not Exist")
}

// UserCreate 创建用户
func (*User) UserCreate(req *service.UserRequest) (user User, err error) {
	var count int64
	DB.Where("user_name=?", req.UserName).Count(&count)
	if count != 0 {
		return User{}, errors.New("UserName Exist")
	}

	user = User{
		UserName: req.UserName,
		NickName: req.NickName,
	}

	// 密码的加密
	_ = user.SetPassword(req.Password)
	err = DB.Create(&user).Error
	return user, err
}

// SetPassword 加密密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		user.PasswordConfirm = string(bytes)
	}

	user.PasswordConfirm = string(bytes) // 存入Password字段

	return nil
}

// CheckPassword 检验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordConfirm), []byte(password))
	return err == nil
}

// BuildUser 序列化User
func BuildUser(item User) *service.UserModel {
	userModel := service.UserModel{
		UserID:   uint32(item.UserID),
		UserName: item.UserName,
		NickName: item.NickName,
	}

	return &userModel
}
