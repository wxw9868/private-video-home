package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
	"github.com/wxw9868/video/service"
)

type ActressAdd struct {
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

// ActressAddApi godoc
//
//	@Summary		添加演员
//	@Description	添加演员
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ActressAdd	true	"Actress Add"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/actress/add [post]
func ActressAddApi(c *gin.Context) {
	var bind ActressAdd
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

type ActressEdit struct {
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

// ActressEditApi godoc
//
//	@Summary		编辑演员信息
//	@Description	编辑演员信息
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ActressEdit	true	"Actress Edit"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/actress/edit [post]
func ActressEditApi(c *gin.Context) {
	var bind ActressEdit
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := actressService.Edit(bind.Id, bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("修改成功", nil))
}

// ActressDeleteApi godoc
//
//	@Summary		删除演员
//	@Description	get string by ID
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Actress ID"
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/actress/delete/{id} [get]
func ActressDeleteApi(c *gin.Context) {
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

// ActressListApi godoc
//
//	@Summary		演员列表
//	@Description	演员列表
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ActressList	true	"演员列表"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/actress/list [post]
func ActressListApi(c *gin.Context) {
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

// ActressInfoApi godoc
//
//	@Summary		演员信息
//	@Description	get string by ID
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Actress ID"
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/actress/info/{id} [get]
func ActressInfoApi(c *gin.Context) {
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

// OneAddInfoToActress godoc
//
//	@Summary		获取演员信息
//	@Description	获取演员信息
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Param			actress	query		string	true	"actress"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/actress/oneAddInfo [get]
func OneAddInfoToActress(c *gin.Context) {
	var actress = c.Query("actress")
	if err := service.OneAddInfoToActress(actress); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}

// AllAddInfoToActress godoc
//
//	@Summary		获取所有演员信息
//	@Description	获取所有演员信息
//	@Tags			actress
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/actress/allAddInfo [get]
func AllAddInfoToActress(c *gin.Context) {
	if err := service.AllAddInfoToActress(); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}
