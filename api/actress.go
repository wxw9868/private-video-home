package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
)

type ActressAdd struct {
	Name         string `form:"name" json:"name" binding:"required"`
	Alias        string `form:"alias" json:"alias"`
	Avatar       string `form:"avatar" json:"avatar"`
	Birth        string `form:"birth" json:"birth"`
	Measurements string `form:"measurements" json:"measurements"`
	CupSize      string `form:"cup_size" json:"cup_size"`
	DebutDate    string `form:"debut_date" json:"debut_date"`
	StarSign     string `form:"star_sign" json:"star_sign"`
	BloodGroup   string `form:"blood_group" json:"blood_group"`
	Stature      string `form:"stature" json:"stature"`
	Nationality  string `form:"nationality" json:"nationality"`
	Introduction string `form:"introduction" json:"introduction"`
}

func ActressAddApi(c *gin.Context) {
	var bind ActressAdd
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := as.Add(bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("添加成功", nil))
}

type ActressEdit struct {
	Id   uint   `json:"id" binding:"required"`
	Name string `form:"name" json:"name" binding:"required"`
}

func ActressEditApi(c *gin.Context) {
	var bind ActressEdit
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := as.Edit(bind.Id, bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("修改成功", nil))
}

type ActressDelete struct {
	ID uint `form:"id" json:"id" binding:"required"`
}

func ActressDeleteApi(c *gin.Context) {
	var bind ActressDelete
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := as.Delete(bind.ID); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("删除成功", nil))
}

type ActressList struct {
	Page   int    `uri:"page" form:"page" json:"page"`
	Size   int    `uri:"size" form:"size" json:"size"`
	Action string `uri:"action" form:"action" json:"action"`
	Sort   string `uri:"sort" form:"sort" json:"sort"`
}

func ActressListApi(c *gin.Context) {
	var bind ActressList
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	actresss, err := as.List(bind.Page, bind.Size, bind.Action, bind.Sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("演员列表", map[string]interface{}{
		"list": actresss,
	}))
}
