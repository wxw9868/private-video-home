package api

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
	"github.com/wxw9868/video/model"
)

type ImportFile struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// ImportExcelApi godoc
//
//	@Summary		导入数据
//	@Schemes        http https
//	@Description	用于导入数据
//	@Tags			stock
//	@Accept			json
//	@Produce		json
//	@Param			file	formData	file	true	"文件"
//	@Success		200		{object}	Message
//	@Failure        400     {object}    Message
//	@Failure        404     {object}    Message
//	@Failure        500     {object}    Message
//	@Router			/stock/importExcel [post]
func ImportTradingRecordsApi(c *gin.Context) {
	var bindFile ImportFile
	if err := c.ShouldBind(&bindFile); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	isExist := func(ext string) bool {
		exts := map[string]struct{}{
			".xlsx": {},
			".xls":  {},
			".csv":  {},
		}
		_, ok := exts[ext]
		return ok
	}

	file := bindFile.File
	filename := file.Filename
	ext := filepath.Ext(filename)
	if !isExist(ext) {
		c.JSON(http.StatusBadRequest, util.Fail("文件格式必须是.xlsx .xls .csv其中之一"))
		return
	}

	// 单文件
	// 上传文件至指定的完整文件路径
	dst := "./assets/file/excel/" + filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(fmt.Sprintf("upload file err: %s", err.Error())))
		return
	}

	if err := stock.ImportTradingRecords(dst); err != nil {
		isFileExist := func(name string) bool {
			if _, err := os.Lstat(name); err != nil {
				return os.IsExist(err)
			}
			return true
		}
		if isFileExist(dst) {
			if err = os.Remove(dst); err != nil {
				fmt.Errorf("文件删除失败: %s", err)
			}
		}
		fmt.Errorf("文件导入失败: %s", err)
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("导入成功", nil))
}

// LiquidationApi godoc
//
//	@Summary		已清仓股票
//	@Schemes        http https
//	@Description	用于已清仓股票
//	@Tags			stock
//	@Accept			json
//	@Produce		json
//	@Param          page      query     string  false    "页码"
//	@Param          page_size query     string  false   "每页条数"
//	@Success		200		  {object}	Message
//	@Failure        400       {object}  Message
//	@Failure        404       {object}  Message
//	@Failure        500       {object}  Message
//	@Router			/stock/liquidation [get]
func LiquidationApi(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	list, err := stock.Liquidation(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	paginator, err := stock.Pagination(&model.Liquidation{}, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
		"paginator": gin.H{
			"totalCount":  paginator.TotalCount(),
			"totalPage":   paginator.TotalPage(),
			"prePage":     paginator.PrePage(),
			"currentPage": paginator.CurrentPage(),
			"nextPage":    paginator.NextPage(),
			"pageRange":   paginator.Pages(),
		},
	})
}

type Paginate struct {
	Page     uint `form:"page" json:"page"`
	PageSize uint `form:"page_size" json:"page_size"`
}

// TradingRecordsApi godoc
//
//	@Summary		历史成交
//	@Schemes        http https
//	@Description	历史成交数据
//	@Tags			stock
//	@Accept			json
//	@Produce		json
//	@Param          page      query     string  false    "页码"
//	@Param          page_size query     string  false    "每页条数"
//	@Success		200		  {object}	Message
//	@Failure        400       {object}  Message
//	@Failure        404       {object}  Message
//	@Failure        500       {object}  Message
//	@Router			/stock/tradingRecords [get]
func TradingRecordsApi(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	list, err := stock.TradingRecords(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	paginator, err := stock.Pagination(&model.TradingRecords{}, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
		"paginator": gin.H{
			"totalCount":  paginator.TotalCount(),
			"totalPage":   paginator.TotalPage(),
			"prePage":     paginator.PrePage(),
			"currentPage": paginator.CurrentPage(),
			"nextPage":    paginator.NextPage(),
			"pageRange":   paginator.Pages(),
		},
	})
}
