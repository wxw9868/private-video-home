package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
	"github.com/wxw9868/video/utils"
)

type Index struct {
	Id       uint32      `json:"id" binding:"required"`
	Text     string      `json:"text" binding:"required"`
	Document interface{} `json:"document" binding:"required"`
}

// SearchIndex godoc
//
//	@Summary		增加/修改索引
//	@Description	增加/修改索引
//	@Tags			search
//	@Accept			json
//	@Produce		json
//	@Param			data	body		Index	true	"增加/修改索引"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/search/api/index [post]
func SearchIndex(c *gin.Context) {
	var bind Index
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	b, err := json.Marshal(&bind)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	Post(c, utils.Join("/index", "?", "database=", "private-video"), bytes.NewReader(b))
}

// SearchIndexBatch godoc
//
//	@Summary		批量增加/修改索引
//	@Description	批量增加/修改索引
//	@Tags			search
//	@Accept			json
//	@Produce		json
//	@Param			data	body		[]Index	true	"批量增加/修改索引"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/search/api/index/batch [post]
func SearchIndexBatch(c *gin.Context) {
	var binds []Index
	if err := c.ShouldBindJSON(&binds); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	b, err := json.Marshal(&binds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	Post(c, utils.Join("/index/batch", "?", "database=", "private-video"), bytes.NewReader(b))
}

type IndexRemove struct {
	Id uint32 `json:"id" binding:"required"`
}

// SearchIndexRemove godoc
//
//	@Summary		删除索引
//	@Description	删除索引
//	@Tags			search
//	@Accept			json
//	@Produce		json
//	@Param			data	body		IndexRemove	true	"删除索引"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/search/api/index/remove [post]
func SearchIndexRemove(c *gin.Context) {
	var bind IndexRemove
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	b, err := json.Marshal(&bind)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	Post(c, utils.Join("/index/remove", "?", "database=", "private-video"), bytes.NewReader(b))
}

type Query struct {
	Query     string      `json:"query" binding:"required"`
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	Order     string      `json:"order"`
	Highlight interface{} `json:"highlight"`
	ScoreExp  string      `json:"scoreExp"`
}

// SearchQuery godoc
//
//	@Summary		查询索引
//	@Description	查询索引
//	@Tags			search
//	@Accept			json
//	@Produce		json
//	@Param			data	body		Query	true	"查询索引"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/search/api/query [post]
func SearchQuery(c *gin.Context) {
	var bind Query
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	b, err := json.Marshal(&bind)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	Post(c, utils.Join("/query", "?", "database=", "private-video"), bytes.NewReader(b))
}

// SearchStatus godoc
//
//	@Summary		查询状态
//	@Description	查询状态
//	@Tags			search
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/search/api/status [get]
func SearchStatus(c *gin.Context) {
	Get(c, "/status")
}

// SearchDbDrop godoc
//
//	@Summary		删除数据库
//	@Description	删除数据库
//	@Tags			search
//	@Accept			json
//	@Produce		json
//	@Param			database	query		string	true	"删除数据库"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/search/api/db/drop [get]
func SearchDbDrop(c *gin.Context) {
	Get(c, utils.Join("/db/drop", "?", "database=", c.Query("database")))
}

// SearchWordCut godoc
//
//	@Summary		在线分词
//	@Description	在线分词
//	@Tags			search
//	@Accept			json
//	@Produce		json
//	@Param			q		query		string	true	"在线分词"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/search/api/word/cut [get]
func SearchWordCut(c *gin.Context) {
	Get(c, utils.Join("/word/cut", "?", "q=", c.Query("q")))
}

func Post(c *gin.Context, url string, body io.Reader) {
	resp, err := client.POST(url, "application/json", body)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	defer resp.Body.Close()

	robots, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	// 对JSON字符串进行格式化（缩进）
	// formattedJSONBytes, err := json.MarshalIndent(json.RawMessage(robots), "", "  ")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
	// 	return
	// }
	// formattedJSONString := string(formattedJSONBytes)

	_, err = c.Writer.Write(robots)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
}

func Get(c *gin.Context, url string) {
	resp, err := client.GET(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	defer resp.Body.Close()

	robots, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	_, err = c.Writer.Write(robots)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
}
