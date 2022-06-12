package v1

import (
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//查询分类名是否存在

//添加分类
func AddCategory(c *gin.Context) {
	var data model.Category
	_ = c.ShouldBindJSON(&data)
	//找到自定义的状态码
	code = model.CheckCategory(data.Name)
	if code == errmsg.SUCCESS {
		model.CreateCate(&data)
	}
	if code == errmsg.ERROR_CATENAME_USED {
		code = errmsg.ERROR_CATENAME_USED
	}
	//http.StatusOK 是网络状态码，不是自定义的
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code), //调用errmsg中的GetErrMsg函数，将自定义的错误信息传递出去
	})
}

//查询单个用户

//查询分类列表
func GetCate(c *gin.Context) {
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

	data,total := model.GetCate(pageSize, pageNum)
	code = errmsg.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":data,
		"total":total,
		"message":errmsg.GetErrMsg(code),
	})
}

//编辑分类
func EditCate(c *gin.Context) {
	//拿到分类ID
	id,_ := strconv.Atoi(c.Param("id"))
	//先查询有没有同名，再编辑用户
	var data model.Category
	c.ShouldBindJSON(&data)
	code = model.CheckCategory(data.Name)
	if code == errmsg.SUCCESS {
		model.EditCate(id,&data) //编辑
	}
	if code == errmsg.ERROR_CATENAME_USED {
		c.Abort()  //直接返回
	}
	c.JSON(http.StatusOK,gin.H{
		"status":code,
		"message":errmsg.GetErrMsg(code),
	})
}

//删除分类
func DeleteCate(c *gin.Context) {
	//拿到分类ID
	id,_ := strconv.Atoi(c.Param("id"))
	code = model.DeleteCate(id)
	c.JSON(http.StatusOK,gin.H{
		"status":code,
		"message":errmsg.GetErrMsg(code),
	})
}

