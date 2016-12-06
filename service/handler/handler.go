package handler

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

const limitUpperBound = 1000

// Get limit and offset from URL parameters
func getLimitOffset(params url.Values) (int, int, error) {
	limit, offset := 10, 0

	limitStr := params.Get("limit")
	if len(limitStr) > 0 {
		customLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, err
		}
		if customLimit > 0 && customLimit <= limitUpperBound {
			limit = customLimit
		}
	}

	offsetStr := params.Get("offset")
	if len(offsetStr) > 0 {
		customOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, err
		}
		if customOffset > 0 {
			offset = customOffset
		}
	}
	return limit, offset, nil
}

// Error response for handlers
func respondWithErr(c *gin.Context, err error, code int, errCode string) {
	log.Println(err)
	c.JSON(code, map[string]interface{}{"success": false, "msg": errCode})
}

// Success response for handlers
func respondWithSuccess(c *gin.Context, code int) {
	c.JSON(code, map[string]interface{}{"success": true})
}

// Success response for search handlers
func respondSearchItems(c *gin.Context, total, offset, limit int, items interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"total":   total,
		"offset":  offset,
		"limit":   limit,
		"items":   items,
	})
}
