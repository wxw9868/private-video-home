package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
)

type TagAdd struct {
	Name string `form:"name" json:"name" binding:"required"`
}

// TagAddApi godoc
//
//	@Summary		Tag Add
//	@Description	Tag Add
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			data	body		TagAdd	true	"Tag Add"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/tag/add [post]
func TagAddApi(c *gin.Context) {
	var bind TagAdd
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := tagService.Add(bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("添加成功", nil))
}

type TagEdit struct {
	Id   uint   `json:"id" binding:"required"`
	Name string `form:"name" json:"name" binding:"required"`
}

// TagEditApi godoc
//
//	@Summary		Tag Edit
//	@Description	Tag Edit
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			data	body		TagEdit	true	"Tag Edit"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/tag/edit [post]
func TagEditApi(c *gin.Context) {
	var bind TagEdit
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := tagService.Edit(bind.Id, bind.Name); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("修改成功", nil))
}

// TagDeleteApi godoc
//
//	@Summary		Tag Delete
//	@Description	Tag Delete
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"Tag Delete"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/tag/delete/{id} [get]
func TagDeleteApi(c *gin.Context) {
	id := c.Param("video_id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	if aid == 0 {
		c.JSON(http.StatusBadRequest, util.Fail("video_id must be greater than 0"))
		return
	}

	if err := tagService.Delete(uint(aid)); err != nil {
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

// TagListApi godoc
//
//	@Summary		Tag List
//	@Description	Tag List
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			data	body		TagList	true	"Tag List"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/tag/list [post]
func TagListApi(c *gin.Context) {
	var bind TagList
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	tags, err := tagService.List(bind.Page, bind.Size, bind.Action, bind.Sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("标签列表", map[string]interface{}{
		"list": tags,
	}))
}
