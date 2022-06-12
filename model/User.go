package model

import (
	"encoding/base64"
	"ginblog/utils/errmsg"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	//validate : 数据验证
	Username string `gorm:"type:varchar(20);not null" json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(20);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	//权限值 1是管理员，2是普通用户
	Role int `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

//查询用户是否存在
func CheckUser(name string) (code int) {
	var users User
	db.Select("id").Where("username=?", name).First(&users)
	if users.ID > 0 {
		return errmsg.ERROR_USERNAME_USED //说明用户存在
	}
	return errmsg.SUCCESS
}

//新增用户
func CreateUser(data *User) int {
	//加密密码
	//data.Password = ScryptPw(data.Password)  //用下面的钩子函数实现
	//或者直接使用
	//data.BeforeSave()
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

//分页
//查询用户列表
//pageSize  分页大小   pageNum 分页数量
func GetUsers(pageSize, pageNum int) ([]User,int64) {
	var users []User
	var total int64
	//(pageNum-1)*pageSize  设置偏移量
	err = db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Count(&total).Error
	//gorm.ErrRecordNotFound  记录没有找到的错误
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil,0
	}
	return users,total
}

//删除用户
func DeleteUser(id int) int {
	var user User
	err = db.Where("id=?",id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

//编辑用户信息
func EditUser(id int, data *User) int {
	var user User
	//因为用struct更新的话，不会更新非0值（role为0），所以用map
	maps := make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role
	err =  db.Model(&user).Where("id=?",id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
//用钩子函数对密码加密
func (u *User) BeforeSave(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	return nil
}

//密码加密
func ScryptPw(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 32, 18, 16, 20, 24, 58, 100}
	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw
}

//登录验证
func CheckLogin(username, password string) int {
	var user User
	db.Where("username=?",username).First(&user)
	if user.ID == 0 {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	if ScryptPw(password) != user.Password {
		return errmsg.ERROR_PASSWORD_WRONG
	}
	if user.Role != 1 {
		return errmsg.ERROR_USER_NO_RIGHT
	}
	return errmsg.SUCCESS
}
