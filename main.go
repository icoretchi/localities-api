package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Locality struct {
	Code            string `json:"code"`
	StatisticalCode string `json:"statisticalCode"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	ParentCode      string `json:"parentCode"`
}

var localities []Locality

func init() {
	localities = make([]Locality, 0)
}

func NewLocalityHandler(c *gin.Context) {
	var locality Locality
	if err := c.ShouldBindJSON(&locality); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	localities = append(localities, locality)
	c.JSON(http.StatusOK, locality)
}

func main() {
	router := gin.Default()
	router.POST("/localities", NewLocalityHandler)
	router.Run()
}
