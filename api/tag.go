package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
)

type TagAdd struct {
	Name string `form:"name" json:"name" binding:"required"`
}

func TagAddApi(c *gin.Context) {
	var bind TagAdd
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := ts.Add(bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("添加成功", nil))
}

type TagEdit struct {
	Id   uint   `json:"id" binding:"required"`
	Name string `form:"name" json:"name" binding:"required"`
}

func TagEditApi(c *gin.Context) {
	var bind TagEdit
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := ts.Edit(bind.Id, bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("修改成功", nil))
}

type TagDelete struct {
	ID uint `form:"id" json:"id" binding:"required"`
}

func TagDeleteApi(c *gin.Context) {
	var bind TagDelete
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := ts.Delete(bind.ID); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("删除成功", nil))
}

type TagList struct {
	Page   int    `uri:"page" form:"page" json:"page"`
	Size   int    `uri:"size" form:"size" json:"size"`
	Action string `uri:"action" form:"action" json:"action"`
	Sort   string `uri:"sort" form:"sort" json:"sort"`
}

func TagListApi(c *gin.Context) {
	var bind TagList
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	tags, err := ts.List(bind.Page, bind.Size, bind.Action, bind.Sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("标签列表", map[string]interface{}{
		"list": tags,
	}))
}
