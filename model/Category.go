package model

import (
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
)

/*
分类
 */
type Category struct {
	ID   uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}

//查询分类是否存在
func CheckCategory(name string) (code int) {
	var cate Category
	db.Select("id").Where("name=?", name).First(&cate)
	if cate.ID > 0 {
		return errmsg.ERROR_CATENAME_USED //说明用户存在
	}
	return errmsg.SUCCESS
}

//新增分类
func CreateCate(data *Category) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

//分页
//查询分类列表
//pageSize  分页大小   pageNum 分页数量
func GetCate(pageSize, pageNum int) ([]Category,int64) {
	var cate []Category
	var total int64
	//(pageNum-1)*pageSize  设置偏移量
	err = db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&cate).Count(&total).Error
	//gorm.ErrRecordNotFound  记录没有找到的错误
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil,0
	}
	return cate,total
}

//删除分类
func DeleteCate(id int) int {
	var cate Category
	err = db.Where("id=?",id).Delete(&cate).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

//编辑分类信息
func EditCate(id int, data *Category) int {
	var cate Category
	//因为用struct更新的话，不会更新非0值（role为0），所以用map
	maps := make(map[string]interface{})
	maps["name"] = data.Name
	err =  db.Model(&cate).Where("id=?",id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}