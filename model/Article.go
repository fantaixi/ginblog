package model

import (
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
)

type Article struct {
	Category Category `gorm:"foreignkey:Cid"`
	gorm.Model
	Title   string `gorm:"type:varchar(100);not null" json:"title"`
	Cid     int    `gorm:"type:int;not null" json:"cid"`
	Desc    string `gorm:"type:varchar(200);not null" json:"desc"`
	Content string `gorm:"type:longtext"  json:"content"`
	Img     string `gorm:"type:varchar(100)" json:"img"`
}

//新增文章
func CreateArt(data *Article) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

//查询分类下的所有文章
func GetCateArt(id int, pageSize, pageNum int) ([]Article, int, int64) {
	var cateArtList []Article
	var total int64
	err := db.Preload("Category").Limit(pageSize).Offset((pageNum-1)*pageSize).Where("cid=?", id).Find(&cateArtList).Count(&total).Error
	if err != nil {
		return nil, errmsg.ERROR_CATE_NOT_EXIST,0
	}
	return cateArtList, errmsg.SUCCESS,total
}

//查询单个文章
func GetArtInfo(id int) (Article, int) {
	var art Article
	err := db.Preload("Category").Where("id=?", id).First(&art).Error
	if err != nil {
		return art, errmsg.ERROR_ART_NOT_EXIST
	}
	return art, errmsg.SUCCESS
}

//分页
//查询文章列表
//pageSize  分页大小   pageNum 分页数量
func GetArt(pageSize, pageNum int) ([]Article, int,int64) {
	var artList []Article
	var total int64
	//(pageNum-1)*pageSize  设置偏移量
	err = db.Preload("Category").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&artList).Count(&total).Error
	//gorm.ErrRecordNotFound  记录没有找到的错误
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errmsg.ERROR,0
	}
	return artList, errmsg.SUCCESS,total
}

//删除文章
func DeleteArt(id int) int {
	var art Article
	err = db.Where("id=?", id).Delete(&art).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

//编辑文章信息
func EditArt(id int, data *Article) int {
	var art Article
	//因为用struct更新的话，不会更新非0值（role为0），所以用map
	maps := make(map[string]interface{})
	maps["title"] = data.Title
	maps["cid"] = data.Cid
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["img"] = data.Img

	err = db.Model(&art).Where("id=?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
