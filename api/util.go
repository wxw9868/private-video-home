package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
)

// OneAddInfoToActress godoc
//
//	@Summary	通过爬虫获取演员信息
//	@Tags		util
//	@Accept		json
//	@Produce	json
//	@Param		actress	query		string	true	"actress"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/actress/oneAddInfo [get]
func OneAddInfoToActress(c *gin.Context) {
	var actress = c.Query("actress")
	if err := utilService.OneAddInfoToActress(actress); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}

// AllAddInfoToActress godoc
//
//	@Summary	通过爬虫获取所有演员信息
//	@Tags		util
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	Success
//	@Failure	400	{object}	Fail
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/actress/allAddInfo [get]
func AllAddInfoToActress(c *gin.Context) {
	if err := utilService.AllAddInfoToActress(); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}

type GetVideo struct {
	Page int    `uri:"page" form:"page" json:"page"`
	Size int    `uri:"size" form:"size" json:"size"`
	Site string `uri:"site" form:"site" json:"site"`
}

func GetVideoPic(c *gin.Context) {
	var bind GetVideo
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := utilService.AllVideoPic(bind.Page, bind.Size, bind.Site); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}

type OneVideo struct {
	Actress string `uri:"actress" form:"actress" json:"actress"`
	Site    string `uri:"site" form:"site" json:"site"`
}

func OneVideoPic(c *gin.Context) {
	var bind OneVideo
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := utilService.OneVideoPic(bind.Actress, bind.Site); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}
