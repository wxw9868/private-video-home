package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
)

type AddActress struct {
	Name         string `form:"name" json:"name" binding:"required" example:""`
	Alias        string `form:"alias" json:"alias" example:""`
	Avatar       string `form:"avatar" json:"avatar" example:"assets/image/avatar/anonymous.png"`
	Birth        string `form:"birth" json:"birth" example:""`
	Measurements string `form:"measurements" json:"measurements" example:""`
	CupSize      string `form:"cup_size" json:"cup_size" example:""`
	DebutDate    string `form:"debut_date" json:"debut_date" example:""`
	StarSign     string `form:"star_sign" json:"star_sign" example:""`
	BloodGroup   string `form:"blood_group" json:"blood_group" example:""`
	Stature      string `form:"stature" json:"stature" example:""`
	Nationality  string `form:"nationality" json:"nationality" example:""`
	Introduction string `form:"introduction" json:"introduction" example:""`
}

// AddActressApi godoc
//
//	@Summary		添加演员
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Param			data	body		AddActress	true	"演员信息"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/actress/addActress [post]
func AddActressApi(c *gin.Context) {
	var bind AddActress
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := actressService.Add(bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("添加成功", nil))
}

type UpdateActress struct {
	Id           uint   `json:"id" binding:"required" example:""`
	Name         string `form:"name" json:"name" binding:"required" example:""`
	Alias        string `form:"alias" json:"alias" example:""`
	Avatar       string `form:"avatar" json:"avatar" example:""`
	Birth        string `form:"birth" json:"birth" example:""`
	Measurements string `form:"measurements" json:"measurements" example:""`
	CupSize      string `form:"cup_size" json:"cup_size" example:""`
	DebutDate    string `form:"debut_date" json:"debut_date" example:""`
	StarSign     string `form:"star_sign" json:"star_sign" example:""`
	BloodGroup   string `form:"blood_group" json:"blood_group" example:""`
	Stature      string `form:"stature" json:"stature" example:""`
	Nationality  string `form:"nationality" json:"nationality" example:""`
	Introduction string `form:"introduction" json:"introduction" example:""`
}

// UpdateActressApi godoc
//
//	@Summary		更新演员信息
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Param			data	body		UpdateActress	true	"演员信息"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/actress/updateActress [post]
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

// DeleteActressApi godoc
//
//	@Summary		删除演员
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"演员ID"
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/actress/deleteActress/{id} [get]
func DeleteActressApi(c *gin.Context) {
	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	if aid == 0 {
		c.JSON(http.StatusBadRequest, util.Fail("id must be greater than 0"))
		return
	}

	if err := actressService.Delete(uint(aid)); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("删除成功", nil))
}

type ActressList struct {
	Page    int    `uri:"page" form:"page" json:"page"`
	Size    int    `uri:"size" form:"size" json:"size"`
	Action  string `uri:"action" form:"action" json:"action" example:"a.CreatedAt"`
	Sort    string `uri:"sort" form:"sort" json:"sort" example:"desc"`
	Actress string `uri:"actress"  form:"actress"  json:"actress" example:""`
}

// GetActressListApi godoc
//
//	@Summary		演员列表
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ActressList	true	"演员列表"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/actress/getActressList [post]
func GetActressListApi(c *gin.Context) {
	var bind ActressList
	if err := c.BindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	data, err := actressService.List(bind.Page, bind.Size, bind.Action, bind.Sort, bind.Actress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("演员列表", data))
}

// GetActressInfoApi godoc
//
//	@Summary		获取演员信息
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"演员ID"
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/actress/getActressInfo/{id} [get]
func GetActressInfoApi(c *gin.Context) {
	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	if aid == 0 {
		c.JSON(http.StatusBadRequest, util.Fail("id must be greater than 0"))
		return
	}

	actress, err := actressService.Info(uint(aid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("演员信息", actress))
}
