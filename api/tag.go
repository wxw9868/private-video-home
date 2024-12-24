package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
	"github.com/wxw9868/video/model/request"
)

type CreateTag struct {
	Name string `form:"name" json:"name" binding:"required"`
}

// CreateTagApi godoc
//
//	@Summary	Create Tag
//	@Tags		tag
//	@Accept		json
//	@Produce	json
//	@Param		data	body		CreateTag	true	"Create Tag"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/tag/createTag [post]
func CreateTagApi(c *gin.Context) {
	var bind CreateTag
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := tagService.Create(bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("添加成功", nil))
}

type TagList struct {
	request.Paginate
	request.OrderBy
}

// GetTagListApi godoc
//
//	@Summary	Get Tag List
//	@Tags		tag
//	@Accept		json
//	@Produce	json
//	@Param		data	body		TagList	true	"Tag List"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/tag/getTagList [post]
func GetTagListApi(c *gin.Context) {
	var bind TagList
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	data, err := tagService.List(bind.Page, bind.Size, bind.Column, bind.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("标签列表", data))
}

// GetTagInfoApi godoc
//
//	@Summary	Get Tag Info
//	@Tags		tag
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Tag ID"
//	@Success	200	{object}	Success
//	@Failure	400	{object}	Fail
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/tag/getTagInfo/{id} [get]
func GetTagInfoApi(c *gin.Context) {
	var bind request.Common
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	data, err := tagService.Info(bind.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("标签列表", data))
}

type UpdateTag struct {
	ID uint `json:"id" binding:"required"`
	CreateTag
}

// UpdateTagApi godoc
//
//	@Summary	Update Tag
//	@Tags		tag
//	@Accept		json
//	@Produce	json
//	@Param		data	body		UpdateTag	true	"Update Tag"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/tag/updateTag [post]
func UpdateTagApi(c *gin.Context) {
	var bind UpdateTag
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := tagService.Update(bind.ID, bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("修改成功", nil))
}

// DeleteTagApi godoc
//
//	@Summary	Delete Tag
//	@Tags		tag
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Tag ID"
//	@Success	200	{object}	Success
//	@Failure	400	{object}	Fail
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/tag/deleteTag/{id} [get]
func DeleteTagApi(c *gin.Context) {
	var bind request.Common
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := tagService.Delete(bind.ID); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("删除成功", nil))
}
