package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
)

type Actress struct {
	Name         string `form:"name" json:"name" binding:"required"`
	Alias        string `form:"alias" json:"alias"`
	Avatar       string `form:"avatar" json:"avatar" example:"assets/image/avatar/anonymous.png"`
	Birth        string `form:"birth" json:"birth"`
	Measurements string `form:"measurements" json:"measurements"`
	CupSize      string `form:"cup_size" json:"cup_size"`
	DebutDate    string `form:"debut_date" json:"debut_date"`
	StarSign     string `form:"star_sign" json:"star_sign"`
	BloodGroup   string `form:"blood_group" json:"blood_group"`
	Stature      string `form:"stature" json:"stature"`
	Nationality  string `form:"nationality" json:"nationality"`
	Intro        string `form:"intro" json:"intro"`
}

// CreateActressApi godoc
//
//	@Summary	添加演员
//	@Tags		actress
//	@Accept		json
//	@Produce	json
//	@Param		data	body		Actress	true	"演员信息"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/actress/createActress [post]
func CreateActressApi(c *gin.Context) {
	var bind Actress
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := actressService.Create(bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("添加成功", nil))
}

type ActressList struct {
	Paginate
	OrderBy
	Actress string `uri:"actress"  form:"actress"  json:"actress"`
}

// GetActressListApi godoc
//
//	@Summary	演员列表
//	@Tags		actress
//	@Accept		application/x-www-form-urlencoded
//	@Produce	json
//	@Param		page	query		int		false	"页码"	default(1)
//	@Param		size	query		int		false	"条数"	default(10)
//	@Param		column	query		string	false	"排序字段"
//	@Param		order	query		string	false	"排序方式"	Enums(desc, asc)
//	@Param		actress	query		string	false	"演员名称"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/actress/getActressList [post]
func GetActressListApi(c *gin.Context) {
	var bind ActressList
	if err := c.BindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	data, err := actressService.List(bind.Page, bind.Size, bind.Column, bind.Order, bind.Actress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("演员列表", data))
}

// GetActressInfoApi godoc
//
//	@Summary	获取演员信息
//	@Tags		actress
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"演员ID"
//	@Success	200	{object}	Success
//	@Failure	400	{object}	Fail
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/actress/getActressInfo/{id} [get]
func GetActressInfoApi(c *gin.Context) {
	var bind Common
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	actress, err := actressService.Info(bind.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("演员信息", actress))
}

// DeleteActressApi godoc
//
//	@Summary	删除演员
//	@Tags		actress
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"演员ID"
//	@Success	200	{object}	Success
//	@Failure	400	{object}	Fail
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/actress/deleteActress/{id} [delete]
func DeleteActressApi(c *gin.Context) {
	var bind Common
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := actressService.Delete(bind.ID); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("删除成功", nil))
}

type UpdateActress struct {
	Id uint `json:"id" binding:"required" `
	Actress
}

// UpdateActressApi godoc
//
//	@Summary	更新演员信息
//	@Tags		actress
//	@Accept		json
//	@Produce	json
//	@Param		data	body		UpdateActress	true	"演员信息"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/actress/updateActress [post]
func UpdateActressApi(c *gin.Context) {
	var bind UpdateActress
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := actressService.Updates(bind.Id, bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("修改成功", nil))
}
