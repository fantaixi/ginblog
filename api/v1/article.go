package v1

import (
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//添加文章
func AddArticle(c *gin.Context) {
	var data model.Article
	_ = c.ShouldBindJSON(&data)
	//找到自定义的状态码
	code = model.CreateArt(&data)

	//http.StatusOK 是网络状态码，不是自定义的
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code), //调用errmsg中的GetErrMsg函数，将自定义的错误信息传递出去
	})
}

//查询分类下的所有文章
func GetCateArt(c *gin.Context) {
	//转换前端传来的参数
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	id,_ := strconv.Atoi(c.Param("id"))
	if pageSize == 0 {
		//用-1取消抵消条件
		pageSize = -1
	}
	if pageNum == 0 {
		pageNum = -1
	}
	data,code,toal := model.GetCateArt(id,pageSize,pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":data,
		"total":toal,
		"message":errmsg.GetErrMsg(code),
	})
}

//查询单个文章信息
func GetArtInfo(c *gin.Context) {
	//拿到ID
	id,_ := strconv.Atoi(c.Param("id"))
	data,code := model.GetArtInfo(id)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":data,
		"message":errmsg.GetErrMsg(code),
	})
}

//查询文章列表
func GetArt(c *gin.Context) {
	//转换前端传来的参数
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	if pageSize == 0 {
		//用-1取消抵消条件
		pageSize = -1
	}
	if pageNum == 0 {
		pageNum = -1
	}

	data,code,total := model.GetArt(pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":data,
		"total":total,
		"message":errmsg.GetErrMsg(code),
	})
}

//编辑文章
func EditArt(c *gin.Context) {
	//拿到ID
	id,_ := strconv.Atoi(c.Param("id"))
	//先查询有没有同名，再编辑用户
	var data model.Article
	c.ShouldBindJSON(&data)
	code = model.EditArt(id,&data)

	c.JSON(http.StatusOK,gin.H{
		"status":code,
		"message":errmsg.GetErrMsg(code),
	})
}

//删除文章
func DeleteArt(c *gin.Context) {
	//拿到分类ID
	id,_ := strconv.Atoi(c.Param("id"))
	code = model.DeleteArt(id)
	c.JSON(http.StatusOK,gin.H{
		"status":code,
		"message":errmsg.GetErrMsg(code),
	})
}


